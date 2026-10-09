import { error } from '@sveltejs/kit';
import { posts, getPostMeta } from '#lib/posts';

export function entries() {
	return posts.map((p) => ({ slug: p.slug }));
}

export function load({ params }: { params: { slug: string } }) {
	const meta = getPostMeta(params.slug);
	if (!meta) throw error(404, 'Post not found');
	return { slug: params.slug, meta };
}
