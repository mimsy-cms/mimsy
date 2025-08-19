<script lang="ts">
	import { Composer, ContentEditable, RichTextPlugin } from 'svelte-lexical';
	import { theme } from 'svelte-lexical/dist/themes/default';
	import Toolbar from './Toolbar.svelte';

	type Props = {
		composer: Composer;
	}

	let { composer = $bindable() }: Props = $props();

	const initialConfig = {
		theme,
		namespace: 'lexical',
		nodes: [],
		onError: (error: Error) => {
			throw error;
		}
	};
</script>

<Composer {initialConfig} bind:this={composer}>
	<div class="editor-shell svelte-lexical !m-0 w-full !max-w-none">
		<Toolbar />
		<div class="editor-container !border-gray-300">
			<div class="editor-scroller">
				<div class="editor">
					<ContentEditable />
				</div>
			</div>
			<RichTextPlugin />
		</div>
	</div>
</Composer>

<style>
	@reference "tailwindcss/theme";

	:global(.toolbar) {
		@apply !border-gray-300;
	}
</style>
