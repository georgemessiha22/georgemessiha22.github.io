import adapter from '@sveltejs/adapter-static';
import { vitePreprocess } from '@sveltejs/vite-plugin-svelte';
import { mdsvex } from 'mdsvex';

/** @type {import('@sveltejs/kit').Config} */
const config = {
	extensions: ['.svelte', '.svx', '.md'],
	preprocess: [
		vitePreprocess(),
		mdsvex({ extensions: ['.svx', '.md'] })
	],
	kit: {
		adapter: adapter({
			pages: '../../gen/site',
			assets: '../../gen/site',
			fallback: undefined,
			precompress: false,
			strict: true
		}),
		// User page (georgemessiha22.github.io) is served from the domain root.
		paths: { base: '' }
	}
};

export default config;
