<script lang="ts">
	import '../app.css';

	let { children } = $props();
</script>

<svelte:head>
	<link rel="preload" as="font" href="/fonts/inter.ttf" type="font/ttf" crossorigin="anonymous" />
</svelte:head>

{@render children()}
