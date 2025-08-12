<script lang="ts">
	import { Composer } from 'svelte-lexical';
	import Editor from './Editor.svelte';

	let composer: Composer;
</script>

<Editor bind:composer />
