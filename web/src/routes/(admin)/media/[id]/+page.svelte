<script lang="ts">
	import { goto } from '$app/navigation';
	import { formatFileSize } from '$lib/utils/file.js';
	import ConfirmDeletionModal from '$lib/components/ConfirmDeletionDialog.svelte';
	import DownloadIcon from '@lucide/svelte/icons/download';
	import ArrowLeftIcon from '@lucide/svelte/icons/arrow-left';
	import Trash2Icon from '@lucide/svelte/icons/trash-2';
	import CircleXIcon from '@lucide/svelte/icons/circle-x';
	import XIcon from '@lucide/svelte/icons/x';

	let { data } = $props();

	let deleteError = $state<string | null>(null);
	let showDeleteModal = $state(false);
	let isDeleting = $state(false);

	// TODO: Move this to a i18n / localization file
	function formatDate(dateString: string): string {
		const date = new Date(dateString);
		return date.toLocaleDateString('en-US', {
			year: 'numeric',
			month: 'long',
			day: 'numeric',
			hour: '2-digit',
			minute: '2-digit'
		});
	}

	async function handleConfirmDelete() {
		isDeleting = true;
		deleteError = null;

		try {
			const response = await fetch(`/api/v1/media/${data.media.id}`, {
				method: 'DELETE',
				headers: {
					'Content-Type': 'application/json'
				}
			});

			if (!response.ok) {
				if (response.status === 409) {
					throw new Error('Media is in use and cannot be deleted');
				} else {
					throw new Error('Failed to delete media');
				}
			}

			goto('/collections/media');
		} catch (error) {
			deleteError = error instanceof Error ? error.message : 'Failed to delete media';
		} finally {
			isDeleting = false;
			showDeleteModal = false;
		}
	}

	function handleCancelDelete() {
		showDeleteModal = false;
		deleteError = null;
	}
</script>

<div class="flex flex-col gap-6">
	{#if deleteError}
		<div class="rounded-md border border-red-200 bg-red-50 p-4">
			<div class="flex flex-row items-center gap-4">
				<div class="flex-shrink-0">
					<CircleXIcon class="h-5 w-5 text-red-400" />
				</div>
				<p class="flex-1 text-sm text-red-700">
					{deleteError}
				</p>
				<div class="flex-shrink-0">
					<button
						type="button"
						class="cursor-pointer text-gray-600 hover:text-gray-900"
						onclick={() => (deleteError = null)}
					>
						<XIcon class="h-4 w-4" />
					</button>
				</div>
			</div>
		</div>
	{/if}

	<div class="flex items-center gap-4">
		<a
			href="/collections/media"
			class="flex items-center gap-2 text-gray-600 transition-colors hover:text-gray-900"
		>
			<ArrowLeftIcon class="size-4" />
			<span>Back to Media</span>
		</a>
	</div>

	<div class="flex flex-col gap-6">
		<div class="space-y-4">
			<div class="flex flex-col items-center justify-between gap-4 md:flex-row">
				<h1 class="break-all text-3xl font-semibold">{data.media.name}</h1>

				<div class="flex justify-end gap-2">
					<a
						href={data.media.url}
						target="_blank"
						rel="noopener noreferrer"
						class="btn btn-outline flex items-center gap-3"
					>
						<DownloadIcon class="size-5 flex-shrink-0" />
						<span>Download</span>
					</a>
					<button
						type="button"
						class="btn btn-outline flex items-center gap-3"
						onclick={() => (showDeleteModal = true)}
					>
						<Trash2Icon class="size-5 flex-shrink-0" />
						<span>Delete</span>
					</button>
				</div>
			</div>
		</div>

		<div class="flex-coll flex flex-col gap-6 md:flex-row">
			<div class="flex-2/3 overflow-hidden rounded-lg border border-gray-200 bg-white">
				<div class="flex aspect-video items-center justify-center bg-gray-100">
					{#if data.media.content_type.startsWith('image/')}
						<img src={data.media.url} alt={data.media.name} />
					{:else}
						<div class="text-center text-gray-500">
							<div class="text-lg font-medium">No Preview for this file media.</div>
							<div class="text-sm">{data.media.content_type}</div>
						</div>
					{/if}
				</div>
			</div>

			<div class="flex-1/3 rounded-lg border border-gray-200 bg-white p-6">
				<h2 class="mb-4 text-xl font-semibold">Details</h2>

				<dl class="space-y-4">
					<div>
						<dt class="text-sm font-medium text-gray-500">Name</dt>
						<dd class="mt-1 text-sm text-gray-900">{data.media.name}</dd>
					</div>

					<div>
						<dt class="text-sm font-medium text-gray-500">File Size</dt>
						<dd class="mt-1 text-sm text-gray-900">{formatFileSize(data.media.size)}</dd>
					</div>

					<div>
						<dt class="text-sm font-medium text-gray-500">Content Type</dt>
						<dd class="mt-1 text-sm text-gray-900">{data.media.content_type}</dd>
					</div>

					<div>
						<dt class="text-sm font-medium text-gray-500">Upload Date</dt>
						<dd class="mt-1 text-sm text-gray-900">{formatDate(data.media.created_at)}</dd>
					</div>

					<div>
						<dt class="text-sm font-medium text-gray-500">UUID</dt>
						<dd class="mt-1 break-all text-sm text-gray-900">{data.media.uuid}</dd>
					</div>
				</dl>
			</div>
		</div>
	</div>
</div>

<ConfirmDeletionModal
	bind:open={showDeleteModal}
	title="Delete Media"
	description="Are you sure you want to delete {data.media.name} ?"
	onConfirm={handleConfirmDelete}
	onCancel={handleCancelDelete}
/>
