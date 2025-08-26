<script lang="ts">
	import { fly } from 'svelte/transition';
	import { quintOut } from 'svelte/easing';

	const {
		message,
		type = 'error'
	}: {
		message: string;
		type?: 'success' | 'error' | 'warning' | 'info';
	} = $props();

	const styles = {
		success: 'border-green-400 bg-green-100 text-green-700',
		error: 'border-red-400 bg-red-100 text-red-700',
		warning: 'border-yellow-400 bg-yellow-100 text-yellow-700',
		info: 'border-blue-400 bg-blue-100 text-blue-700'
	};

	const icons = {
		success: '✓',
		error: '✕',
		warning: '⚠',
		info: 'ℹ'
	};
</script>

<div 
	class="mb-4 flex items-center gap-3 rounded border px-4 py-3 {styles[type]}"
	transition:fly={{ y: -20, duration: 300, easing: quintOut }}
>
	<span class="text-lg font-bold">{icons[type]}</span>
	<span>{message}</span>
</div>