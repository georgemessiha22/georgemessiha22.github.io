<script lang="ts">
	import type { Download } from '#lib/types';
	import * as DropdownMenu from '#lib/components/ui/dropdown-menu';
	import { Download as DownloadIcon, ChevronDown } from '@lucide/svelte';
	import { cn } from '#lib/utils';

	interface Props {
		downloads: Download[];
	}
	let { downloads }: Props = $props();

	const groups: Array<Download['group']> = ['Summary', 'Detailed'];
</script>

<DropdownMenu.Root>
	<DropdownMenu.Trigger
		class="accent-gradient inline-flex h-8 items-center gap-2 rounded-lg px-3 text-sm font-semibold shadow-sm transition-opacity hover:opacity-90 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
	>
		<DownloadIcon class="size-4" />
		Download Résumé
		<ChevronDown class="size-3.5" />
	</DropdownMenu.Trigger>
	<DropdownMenu.Content align="end" class="w-72">
		{#each groups as g}
			{@const items = downloads.filter((d) => d.group === g)}
			{#if items.length}
				<DropdownMenu.Group>
					<DropdownMenu.GroupHeading class="text-muted-foreground">{g}</DropdownMenu.GroupHeading>
					{#each items as d}
						<DropdownMenu.Item>
							{#snippet child({ props })}
								<a
									href={d.href}
									target="_blank"
									rel="noopener"
									{...props}
									class={cn(props.class as string | undefined, 'flex items-center justify-between')}
								>
									<span>{d.label}</span>
									<span
										class="rounded border border-border px-1.5 py-0.5 text-[10px] text-muted-foreground"
										>{d.kind}</span
									>
								</a>
							{/snippet}
						</DropdownMenu.Item>
					{/each}
				</DropdownMenu.Group>
			{/if}
		{/each}
	</DropdownMenu.Content>
</DropdownMenu.Root>
