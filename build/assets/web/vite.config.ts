import adapter from '@sveltejs/adapter-static';
import { sveltekit } from '@sveltejs/kit/vite';
import { vitePreprocess } from '@sveltejs/vite-plugin-svelte';
import tailwindcss from '@tailwindcss/vite';
import { mdsvex } from 'mdsvex';
import { defineConfig } from 'vite';

export default defineConfig({
	plugins: [
		tailwindcss(),
		sveltekit({
			extensions: ['.svelte', '.svx', '.md'],
			preprocess: [vitePreprocess(), mdsvex({ extensions: ['.svx', '.md'] })],
			adapter: adapter({
				pages: '../../gen/site',
				assets: '../../gen/site',
				fallback: undefined,
				precompress: false,
				strict: true
			}),
			// User page (georgemessiha22.github.io) is served from the domain root.
			paths: { base: '' }
		})
	]
});
