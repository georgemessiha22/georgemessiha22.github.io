// Prepares generated content for the SvelteKit build:
//   1. Runs the Go generator to emit src/lib/data/resume.json (+ static/profile.jpg).
//   2. Syncs repo-root blog/<slug>/index.md posts into src/lib/posts/<slug>/index.md
//      and their images into static/blog/<slug>/ so mdsvex can render them.
//   3. Writes src/lib/data/posts.json (card metadata) derived from each post's
//      front-matter, falling back to the first heading / paragraph / slug.
//
// Safe to run repeatedly; generated directories are recreated each time.
import { spawnSync } from 'node:child_process';
import { existsSync } from 'node:fs';
import fs from 'node:fs/promises';
import path from 'node:path';
import { fileURLToPath } from 'node:url';

const webDir = path.resolve(path.dirname(fileURLToPath(import.meta.url)), '..');
const repoRoot = path.resolve(webDir, '../../..');

const dataFile = path.join(webDir, 'src', 'lib', 'data', 'resume.json');
const postsMeta = path.join(webDir, 'src', 'lib', 'data', 'posts.json');
const postsDir = path.join(webDir, 'src', 'lib', 'posts');
const assetsDir = path.join(webDir, 'static', 'blog');
const blogDir = path.join(repoRoot, 'build', 'assets', 'blog');

async function generateResumeData() {
	const res = spawnSync('go', ['run', './cmd/resume', 'web', '--out', 'build/assets/web'], {
		cwd: repoRoot,
		stdio: 'inherit'
	});
	if (res.status === 0) {
		console.log('[prepare] resume.json generated via Go');
		return;
	}
	if (existsSync(dataFile)) {
		console.warn('[prepare] Go generation failed; using existing resume.json');
		return;
	}
	console.error(
		'[prepare] ERROR: could not generate resume.json (Go failed and no cached file exists).'
	);
	process.exit(1);
}

async function rmrf(dir) {
	await fs.rm(dir, { recursive: true, force: true });
}

// Split leading `---` YAML front-matter from the markdown body.
function splitFrontMatter(src) {
	const m = /^---\r?\n([\s\S]*?)\r?\n---\r?\n?/.exec(src);
	if (!m) return { fm: /** @type {Record<string,string>} */ ({}), body: src };
	/** @type {Record<string,string>} */
	const fm = {};
	for (const line of m[1].split(/\r?\n/)) {
		const kv = /^([A-Za-z0-9_-]+)\s*:\s*(.*)$/.exec(line);
		if (!kv) continue;
		let val = kv[2].trim();
		if (
			(val.startsWith('"') && val.endsWith('"')) ||
			(val.startsWith("'") && val.endsWith("'"))
		) {
			val = val.slice(1, -1);
		}
		fm[kv[1].toLowerCase()] = val;
	}
	return { fm, body: src.slice(m[0].length) };
}

function humanize(slug) {
	return slug
		.replace(/[-_]+/g, ' ')
		.replace(/\b\w/g, (c) => c.toUpperCase())
		.trim();
}

function firstHeading(body) {
	const m = /^#{1,3}\s+(.+?)\s*$/m.exec(body);
	return m ? m[1].trim() : '';
}

function firstParagraph(body) {
	const lines = body.split(/\r?\n/);
	let inFence = false;
	for (const raw of lines) {
		const line = raw.trim();
		if (line.startsWith('```')) {
			inFence = !inFence;
			continue;
		}
		if (inFence || !line) continue;
		if (/^([#>|*_-]|\d+\.|!\[|<)/.test(line)) continue;
		// Strip basic inline markdown for a clean summary.
		const text = line
			.replace(/\*\*(.+?)\*\*/g, '$1')
			.replace(/`([^`]+)`/g, '$1')
			.replace(/\[(.+?)\]\((.+?)\)/g, '$1');
		return text.length > 180 ? text.slice(0, 177).trimEnd() + '…' : text;
	}
	return '';
}

async function syncBlog() {
	await rmrf(postsDir);
	await rmrf(assetsDir);
	await fs.mkdir(postsDir, { recursive: true });
	await fs.mkdir(path.dirname(postsMeta), { recursive: true });

	const metas = [];

	if (existsSync(blogDir)) {
		const entries = await fs.readdir(blogDir, { withFileTypes: true });
		for (const entry of entries) {
			if (!entry.isDirectory()) continue;
			const slug = entry.name;
			const srcPost = path.join(blogDir, slug, 'index.md');
			if (!existsSync(srcPost)) continue;

			const src = await fs.readFile(srcPost, 'utf8');
			const { fm, body } = splitFrontMatter(src);

			metas.push({
				slug,
				title: fm.title || firstHeading(body) || humanize(slug),
				date: fm.date || '',
				summary: fm.summary || firstParagraph(body) || ''
			});

			// Copy the markdown post for mdsvex.
			const destPost = path.join(postsDir, slug, 'index.md');
			await fs.mkdir(path.dirname(destPost), { recursive: true });
			await fs.copyFile(srcPost, destPost);

			// Copy sibling assets (images, etc.) to static/blog/<slug>/.
			const files = await fs.readdir(path.join(blogDir, slug), { withFileTypes: true });
			for (const f of files) {
				if (!f.isFile() || f.name === 'index.md' || f.name.startsWith('.')) continue;
				const destAsset = path.join(assetsDir, slug, f.name);
				await fs.mkdir(path.dirname(destAsset), { recursive: true });
				await fs.copyFile(path.join(blogDir, slug, f.name), destAsset);
			}
		}
	} else {
		console.log('[prepare] no blog/ directory; skipping posts');
	}

	// Newest first; undated posts sort last, ties broken by slug.
	metas.sort((a, b) => {
		if (!!a.date !== !!b.date) return a.date ? -1 : 1;
		if (a.date !== b.date) return a.date < b.date ? 1 : -1;
		return a.slug < b.slug ? -1 : 1;
	});

	await fs.writeFile(postsMeta, JSON.stringify(metas, null, 2) + '\n');
	console.log(`[prepare] synced ${metas.length} blog post(s)`);
}

await generateResumeData();
await syncBlog();
