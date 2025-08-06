<script lang="ts">
	import { cn } from '$lib/cn';
	import Dropzone from '$lib/components/admin/Dropzone.svelte';
	import MediaCard from '$lib/components/admin/media/MediaCard.svelte';
	import UploadProgressPopup from '$lib/components/admin/media/UploadProgressPopup.svelte';
	import { uploadFile, createUploadProgress, type UploadProgress } from '$lib/utils/upload';
	import CloudUploadIcon from '@lucide/svelte/icons/cloud-upload';
	import GridIcon from '@lucide/svelte/icons/grid-3x3';
	import ListIcon from '@lucide/svelte/icons/list';

	let fileInputElement = $state<HTMLInputElement>();
	let layoutMode = $state<'grid' | 'list'>('grid');
	let uploads = $state<UploadProgress[]>([]); // TODO: We might want to use a reactive Set/Map here

	let { data } = $props();

	async function handleFileDrop(files: FileList) {
		const newUploads = createUploadProgress(files);
		uploads.push(...newUploads);

		const uploadPromises = newUploads.map(async (uploadItem) => {
			try {
				const formData = new FormData();
				formData.append('file', uploadItem.file);

				await uploadFile(formData, uploadItem.id, {
					url: '/api/v1/collections/media',
					onProgress: (uploadId, progress) => {
						uploads = uploads.map((u) => (u.id === uploadId ? { ...u, progress } : u));
					},
					onStatusChange: (uploadId, status, error) => {
						uploads = uploads.map((u) => {
							if (u.id === uploadId) {
								const updated = { ...u, status, error };
								// We make sure the progress is set to 100% when
								// the file upload updates to completed.
								if (status === 'completed') {
									updated.progress = 100;
								}
								return updated;
							}
							return u;
						});
					}
				});
			} catch (error) {
				console.error(`Failed to upload ${uploadItem.file.name}:`, error);
			}
		});

		await Promise.all(uploadPromises);
	}

	function clearUploads() {
		uploads = [];
	}
</script>

<div class="flex flex-col gap-6">
	<div class="flex flex-col">
		<h1 class="text-4xl font-medium">Media</h1>

		<div class="mt-4 flex items-center justify-between">
			<div class="flex items-center gap-2">
				<div class="flex overflow-hidden rounded-md border border-gray-300">
					<button
						class={cn('flex items-center gap-2 px-3 py-1.5 text-sm transition-colors', {
							'bg-blue-700 text-white': layoutMode === 'grid',
							'bg-white text-gray-700 hover:bg-gray-50': layoutMode !== 'grid'
						})}
						onclick={() => (layoutMode = 'grid')}
					>
						<GridIcon class="size-4" />
						<span>Grid</span>
					</button>
					<button
						class={cn(
							'flex items-center gap-2 border-l border-gray-300 px-3 py-1.5 text-sm transition-colors',
							{
								'bg-blue-700 text-white': layoutMode === 'list',
								'bg-white text-gray-700 hover:bg-gray-50': layoutMode !== 'list'
							}
						)}
						onclick={() => (layoutMode = 'list')}
					>
						<ListIcon class="size-4" />
						<span>List</span>
					</button>
				</div>
			</div>

			<button class="btn btn-sm flex items-center gap-2" onclick={() => fileInputElement?.click()}>
				<CloudUploadIcon class="size-4" />
				<span>Upload</span>
			</button>
		</div>
	</div>

	<Dropzone id="dropzone" name="dropzone" onChange={handleFileDrop}>
		{#if layoutMode === 'grid'}
			<div
				class="grid grid-cols-1 gap-4 sm:grid-cols-2 md:grid-cols-4 lg:grid-cols-5 xl:grid-cols-6"
			>
				{#each data.media as media}
					<MediaCard
						href={`/collections/media/${media.id}`}
						url={media.url}
						alt={media.name}
						class="transition-transform duration-75 hover:scale-105"
					/>
				{/each}
			</div>
		{:else}
			<div class="w-full overflow-hidden rounded-md border border-gray-200 bg-white">
				<table class="w-full divide-y divide-gray-200">
					<thead class="bg-gray-50">
						<tr>
							<th
								scope="col"
								class="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500"
							>
								Preview
							</th>
							<th
								scope="col"
								class="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500"
							>
								Name
							</th>
							<th
								scope="col"
								class="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500"
							>
								Size
							</th>
							<th
								scope="col"
								class="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500"
							>
								Updated
							</th>
						</tr>
					</thead>
					<tbody class="divide-y divide-gray-200">
						{#each data.media as media}
							<tr class="text-left hover:bg-gray-50">
								<td class="px-6 py-3">
									<div class="h-12 w-12 overflow-hidden rounded-md bg-gray-200">
										<img src={media.url} alt={media.alt} class="h-full w-full object-cover" />
									</div>
								</td>
								<td class="px-6 py-3">
									<a
										class="text-sm text-gray-900 hover:text-blue-600"
										href={`/collections/media/${media.id}`}
									>
										{media.name}
									</a>
								</td>
								<td class="whitespace-nowrap px-6 py-3">
									<div class="text-sm text-gray-500">
										{media.size}
									</div>
								</td>
								<td class="whitespace-nowrap px-6 py-3">
									<div class="text-sm text-gray-500">
										{media.updatedAt}
									</div>
								</td>
							</tr>
						{/each}
					</tbody>
				</table>
			</div>
		{/if}
	</Dropzone>
</div>

<input
	bind:this={fileInputElement}
	type="file"
	multiple
	class="hidden"
	aria-hidden="true"
	onchange={(e) => {
		if (e.currentTarget?.files) {
			handleFileDrop(e.currentTarget?.files);
		}
	}}
/>

<UploadProgressPopup {uploads} onClose={clearUploads} />
