<script lang="ts">
	import Searchbar from '$lib/components/admin/Searchbar.svelte';
	import MoreVerticalIcon from '@lucide/svelte/icons/more-vertical';

	let { data } = $props();
</script>

<div class="flex flex-col gap-6">
	<h1 class="text-4xl font-medium">
		{data.collectionName}
	</h1>

	<Searchbar id="search" name="search" class="max-w-md" />

	<div class="w-full overflow-hidden rounded-md border border-gray-200 bg-white">
		<table class="w-full table-fixed divide-y divide-gray-200">
			<thead class="bg-gray-50">
				<tr>
					<th
						scope="col"
						class="px-6 py-3 text-left text-xs font-medium tracking-wider text-gray-500 uppercase"
					>
						ID
					</th>
					<th
						scope="col"
						class="px-6 py-3 text-left text-xs font-medium tracking-wider text-gray-500 uppercase"
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
								class="flex items-center text-sm text-gray-500"
								href={`/collections/${data.collectionSlug}/${resource.slug}`}
							>
								{resource.slug}
							</a>
						</td>
						<td class="px-6 py-3 text-right text-sm font-medium whitespace-nowrap">
							<div class="flex items-center justify-end space-x-2">
								<button class="text-gray-400 hover:text-gray-600">
									<MoreVerticalIcon class="size-4" />
								</button>
							</div>
						</td>
					</tr>
				{/each}
			</tbody>
		</table>
	</div>
</div>
