<script lang="ts">
	import { cn } from '$lib/cn';
	import MenuIcon from '@lucide/svelte/icons/menu';

	type Props = {
		class?: string;
		onToggle?: () => void;
	};

	let { class: className, onToggle }: Props = $props();
</script>

<div class={cn('border-b border-gray-200 bg-white px-3 py-2', className)}>
	<button class="flex rounded-full p-2 transition-colors hover:bg-neutral-200" onclick={onToggle}>
		<MenuIcon class="size-5 text-gray-600 " />
	</button>
</div>
