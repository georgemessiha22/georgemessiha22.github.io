import type { Component } from 'svelte';
import postsMeta from '#lib/data/posts.json';

export interface PostMeta {
	slug: string;
	title: string;
	date: string;
	summary: string;
}

// Card metadata is produced by scripts/prepare.mjs (posts.json), already sorted
// newest-first. The compiled post components come from mdsvex via import.meta.glob.
export const posts: PostMeta[] = postsMeta as PostMeta[];

interface PostModule {
	default: Component;
}

const modules = import.meta.glob('./posts/*/index.md', { eager: true }) as Record<
	string,
	PostModule
>;

function slugFromPath(p: string): string {
	const m = p.match(/\/posts\/([^/]+)\/index\.md$/);
	return m ? m[1] : p;
}

export function getPostComponent(slug: string): Component | undefined {
	for (const [path, mod] of Object.entries(modules)) {
		if (slugFromPath(path) === slug) return mod.default;
	}
	return undefined;
}

export function getPostMeta(slug: string): PostMeta | undefined {
	return posts.find((p) => p.slug === slug);
}
