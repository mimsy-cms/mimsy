<script lang="ts">
	import { formatFileSize } from '$lib/utils/file';
	import CheckIcon from '@lucide/svelte/icons/check';
	import XIcon from '@lucide/svelte/icons/x';
	import FileIcon from '@lucide/svelte/icons/file';
	import AlertCircleIcon from '@lucide/svelte/icons/alert-circle';
	import type { UploadProgress } from '$lib/utils/upload';
	import { cn } from '$lib/cn';

	type Props = {
		uploads: UploadProgress[];
		onClose?: () => void;
	};

	let { uploads, onClose }: Props = $props();

	const completedUploads = $derived(uploads.filter((u) => u.status === 'completed').length);
	const totalUploads = $derived(uploads.length);
	const hasErrors = $derived(uploads.some((u) => u.status === 'error'));
</script>

{#if uploads.length > 0}
	<div
		class="fixed bottom-4 right-4 z-50 w-96 rounded-lg border border-gray-200 bg-white shadow-lg"
	>
		<div class="flex items-center justify-between border-b border-gray-200 p-4">
			<div class="flex items-center gap-2">
				<h3 class="text-sm font-medium">
					Uploading files ({completedUploads}/{totalUploads})
				</h3>
				{#if hasErrors}
					<span class="text-xs text-red-600">Failed to upload some files</span>
				{/if}
			</div>
			{#if onClose}
				<button
					onclick={onClose}
					class="rounded p-1 text-gray-400 hover:bg-gray-100 hover:text-gray-600"
					aria-label="Close upload progress"
				>
					<XIcon class="size-4" />
				</button>
			{/if}
		</div>

		<div class="max-h-64 overflow-y-auto">
			{#each uploads as upload (upload.id)}
				<div class="border-t border-gray-100 p-4">
					<div class="flex items-center gap-3">
						<div
							class={cn('flex-shrink-0', {
								'text-green-600': upload.status === 'completed',
								'text-red-600': upload.status === 'error',
								'text-blue-600': upload.status === 'uploading'
							})}
						>
							{#if upload.status === 'completed'}
								<CheckIcon class="size-4" />
							{:else if upload.status === 'error'}
								<AlertCircleIcon class="size-4" />
							{:else}
								<FileIcon class="size-4" />
							{/if}
						</div>
						<div class="min-w-0 flex-1">
							<div class="flex items-center justify-between">
								<p class="truncate text-sm font-medium text-gray-900">
									{upload.file.name}
								</p>
								<span class="ml-2 text-xs text-gray-500">
									{formatFileSize(upload.file.size)}
								</span>
							</div>
							{#if upload.status === 'error' && upload.error}
								<p class="mt-1 text-xs text-red-600">{upload.error}</p>
							{:else if upload.status === 'uploading'}
								<div class="mt-2">
									<div class="flex justify-between text-xs text-gray-600">
										<span>Uploading...</span>
										<span>{Math.round(upload.progress)}%</span>
									</div>
									<div class="mt-1 h-1 w-full overflow-hidden rounded-full bg-gray-200">
										<div
											class="h-full bg-blue-600 transition-all duration-300 ease-out"
											style="width: {upload.progress}%"
										></div>
									</div>
								</div>
							{:else if upload.status === 'completed'}
								<p class="mt-1 text-xs text-green-600">Upload completed</p>
							{/if}
						</div>
					</div>
				</div>
			{/each}
		</div>
	</div>
{/if}
