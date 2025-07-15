<script lang="ts">
	import type { Snippet } from 'svelte';

	type Props = {
		children?: Snippet;
	};

	let { children }: Props = $props();
</script>

{#if children}
	<p class="text-xs text-red-500">
		{@render children()}
	</p>
{/if}
