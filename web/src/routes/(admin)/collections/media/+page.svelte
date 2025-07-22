<script lang="ts">
	import Dropzone from '$lib/components/admin/Dropzone.svelte';

	function handleFileDrop(files: FileList) {
		// TODO: Upload files to the API
		console.log(files);
	}
</script>

<div class="flex flex-col gap-6">
	<h1 class="text-4xl font-medium">Media</h1>

	<Dropzone id="dropzone" name="dropzone" onChange={handleFileDrop}>
		<div>
			<img src="https://www.placekittens.com/400/300" alt="" />
			<img src="https://www.placekittens.com/400/300" alt="" />
			<img src="https://www.placekittens.com/400/300" alt="" />
		</div>
	</Dropzone>
</div>
