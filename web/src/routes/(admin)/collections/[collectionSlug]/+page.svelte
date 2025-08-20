<script lang="ts">
	import Searchbar from '$lib/components/admin/Searchbar.svelte';
	import Plus from '@lucide/svelte/icons/plus';

	let { data } = $props();

	function formatDate(dateStr: string) {
		const date = new Date(dateStr);
		const pad = (n: number) => n.toString().padStart(2, '0');

		const day = pad(date.getDate());
		const month = pad(date.getMonth() + 1);
		const year = date.getFullYear();
		const hours = pad(date.getHours());
		const minutes = pad(date.getMinutes());
		const seconds = pad(date.getSeconds());

		return `${day}/${month}/${year}, ${hours}:${minutes}:${seconds}`;
	}
</script>

<div class="flex flex-col gap-6">
	<h1 class="text-4xl font-medium">
		{data.collectionName}
	</h1>

	<div class="flex items-center justify-between">
		<Searchbar id="search" name="search" class="max-w-md" />

		<a
			href={`/collections/${data.collectionSlug}/create`}
			class="inline-flex items-center gap-2 rounded-md bg-blue-600 px-4 py-2 text-white transition-colors hover:bg-blue-700 focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 focus:outline-none"
		>
			<Plus size={16} />
			Create
		</a>
	</div>

	<div class="w-full overflow-hidden rounded-md border border-gray-200 bg-white">
		<table class="w-full table-fixed divide-y divide-gray-200">
			<thead class="bg-gray-50">
				<tr>
					<th
						class="w-1/4 px-6 py-3 text-left text-xs font-medium tracking-wider text-gray-500 uppercase"
					>
						ID
					</th>
					<th
						class="w-1/4 px-6 py-3 text-left text-xs font-medium tracking-wider text-gray-500 uppercase"
					>
						Slug
					</th>
					<th
						class="w-1/4 px-6 py-3 text-left text-xs font-medium tracking-wider text-gray-500 uppercase"
					>
						Created By
					</th>
					<th
						class="w-1/4 px-6 py-3 text-left text-xs font-medium tracking-wider text-gray-500 uppercase"
					>
						Last updated
					</th>
				</tr>
			</thead>
			<tbody class="divide-y divide-gray-200">
				{#each data.resources as resource (resource.id)}
					<tr class="hover:bg-gray-50">
						<td class="px-6 py-3 whitespace-nowrap">
							<div class="flex items-center text-sm text-gray-500">
								{resource.id}
							</div>
						</td>
						<td class="px-6 py-3">
							<a
								class="flex items-center text-sm text-gray-500 hover:text-blue-600 hover:underline"
								href={`/collections/${data.collectionSlug}/${resource.slug}`}
							>
								{resource.slug}
							</a>
						</td>
						<td class="px-6 py-3 whitespace-nowrap">
							<div class="flex items-center text-sm text-gray-500">
								{resource.created_by_email}
							</div>
						</td>
						<td class="px-6 py-3 whitespace-nowrap">
							<div class="flex items-center text-sm text-gray-500">
								{formatDate(resource.updated_at)}
							</div>
						</td>
					</tr>
				{:else}
					<tr>
						<td colspan="4" class="px-6 py-8 text-center text-gray-500">
							<div class="flex flex-col items-center gap-2">
								<p>No resources found</p>
								<a
									href={`/collections/${data.collectionSlug}/create`}
									class="text-blue-600 hover:text-blue-800 hover:underline"
								>
									Create your first resource
								</a>
							</div>
						</td>
					</tr>
				{/each}
			</tbody>
		</table>
	</div>
</div>
