<script lang="ts">
	import Dropzone from '$lib/components/admin/Dropzone.svelte';
	import MediaCard from '$lib/components/admin/media/MediaCard.svelte';
	import CloudUploadIcon from '@lucide/svelte/icons/cloud-upload';

	let fileInputElement = $state<HTMLInputElement>();

	function handleFileDrop(files: FileList) {
		// TODO: Upload files to the API
		console.log(files);
	}
</script>

<div class="flex flex-col gap-6">
	<div class="flex items-end justify-between">
		<h1 class="text-4xl font-medium">Media</h1>

		<button class="btn btn-sm flex items-center gap-2" onclick={() => fileInputElement?.click()}>
			<CloudUploadIcon class="size-4" />
			<span>Upload</span>
		</button>
	</div>

	<Dropzone id="dropzone" name="dropzone" onChange={handleFileDrop}>
		<div class="grid grid-cols-1 gap-4 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5">
			{#each Array(10) as _, index}
				<MediaCard
					href={`/collections/media/${index}`}
					url="https://placehold.co/600x400"
					alt=""
					class="transition-transform duration-75 hover:scale-105"
				/>
			{/each}
		</div>
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
