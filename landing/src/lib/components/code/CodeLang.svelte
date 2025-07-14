<script lang="ts">
	import { highlightCode, normalizeCode } from './prism';

	let { lang, code = '', children = () => null } = $props();

	let normalizedCode = $derived(code ? normalizeCode(code) : '');
	let highlighted = $derived(normalizedCode ? highlightCode(normalizedCode, lang) : '');
</script>

<div data-lang={lang} class="code-lang">
	{#if highlighted}
		<!-- eslint-disable svelte/no-at-html-tags -->
		{@html highlighted}
	{:else}
		{@render children?.()}
	{/if}
</div>

<style>
	.code-lang {
		line-height: 1.4;
		margin: 0;
		padding: 0;
		white-space: pre-wrap;
		display: block;
	}

	.code-lang :global(pre) {
		margin: 0;
		padding: 0;
		line-height: 1.4;
		white-space: pre-wrap;
	}

	.code-lang :global(code) {
		line-height: 1.4;
		white-space: pre-wrap;
	}

	.code-lang :global(p) {
		margin: 0;
		padding: 0;
		line-height: 1.4;
	}

	.code-lang :global(br) {
		line-height: 1.4;
	}
</style>
