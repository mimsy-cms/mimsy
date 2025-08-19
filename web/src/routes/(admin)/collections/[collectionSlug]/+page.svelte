<script lang="ts">
	import Searchbar from '$lib/components/admin/Searchbar.svelte';
	import MoreVerticalIcon from '@lucide/svelte/icons/more-vertical';

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

	<Searchbar id="search" name="search" class="max-w-md" />

	<div class="w-full overflow-hidden rounded-md border border-gray-200 bg-white">
		<table class="table-fixed w-full divide-y divide-gray-200">
			<thead class="bg-gray-50">
				<tr>
					<th class="w-1/4 px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500">
						ID
					</th>
					<th class="w-1/4 px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500">
						Slug
					</th>
					<th class="w-1/4 px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500">
						Created By
					</th>
					<th class="w-1/4 px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500">
						Last updated
					</th>
				</tr>
			</thead>
			<tbody class="divide-y divide-gray-200">
				{#each data.resources as resource}
					<tr class="hover:bg-gray-50">
						<td class="whitespace-nowrap px-6 py-3">
							<div class="flex items-center text-sm text-gray-500">
								{resource.id}
							</div>
						</td>

						<td class="px-6 py-3">
							<a
								class="flex items-center text-sm text-gray-500"
								href={`/collections/${data.collectionSlug}/${resource.slug}`}
							>
								{resource.slug}
							</a>
						</td>

						<td class="whitespace-nowrap px-6 py-3">
							<div class="flex items-center text-sm text-gray-500">
								{resource.created_by_email}
							</div>
						</td>

						<td class="whitespace-nowrap px-6 py-3">
							<div class="flex items-center text-sm text-gray-500">
								{formatDate(resource.updated_at)}
							</div>
						</td>
					</tr>
				{/each}
			</tbody>
		</table>
	</div>
</div>
