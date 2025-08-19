<script lang="ts">
	import { cn } from '$lib/cn';
	import CloudUploadIcon from '@lucide/svelte/icons/cloud-upload';
	import type { Snippet } from 'svelte';
	import { fade } from 'svelte/transition';

	type Props = {
		id: string;
		name: string;
		accept?: string;
		multiple?: boolean;
		disabled?: boolean;
		class?: string;
		children: Snippet;
		onChange: (files: FileList) => void;
	};

	let {
		multiple = false,
		disabled = false,
		class: className,
		onChange,
		children
	}: Props = $props();

	let isDragging = $state(false);
	let containerElement: HTMLDivElement;

	function handleDragOver(event: DragEvent) {
		event.preventDefault();
		if (!disabled) {
			const hasFiles = event.dataTransfer?.types.includes('Files');
			isDragging = hasFiles || false;
		}
	}

	function handleDragLeave(event: DragEvent) {
		event.preventDefault();

		if (event.relatedTarget && containerElement.contains(event.relatedTarget as Node)) {
			return;
		}

		isDragging = false;
	}

	function handleDrop(event: DragEvent) {
		event.preventDefault();
		isDragging = false;

		if (disabled) {
			return;
		}

		const droppedFiles = event.dataTransfer?.files;
		if (droppedFiles) {
			handleFileSelection(droppedFiles);
		}
	}

	function handleFileSelection(files: FileList) {
		if (multiple) {
			const dataTransfer = new DataTransfer();
			Array.from(files).forEach((file) => dataTransfer.items.add(file));
			onChange(dataTransfer.files);
		} else {
			onChange(files);
		}
	}
</script>

<div class={cn('space-y-3', className)} bind:this={containerElement}>
	<div
		class={cn(
			'relative rounded-md border-2 border-dashed border-transparent text-center transition-all duration-200 outline-none',
			isDragging && 'border-blue-500 bg-blue-50',
			disabled && 'cursor-not-allowed opacity-50'
		)}
		ondragover={handleDragOver}
		ondragleave={handleDragLeave}
		ondrop={handleDrop}
		role="button"
		aria-label="Drop files here"
		tabindex={-1}
	>
		{@render children()}

		{#if isDragging}
			<div
				class="absolute inset-0 flex flex-col items-center justify-center rounded-md bg-blue-50/90 text-blue-700 backdrop-blur-sm"
				in:fade={{ duration: 200 }}
				out:fade={{ duration: 150 }}
			>
				<CloudUploadIcon class="size-8" />
				<p class="text-lg font-medium">Drop files to upload</p>
			</div>
		{/if}
	</div>
</div>
