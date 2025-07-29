<script lang="ts">
	import type { PageData } from './$types';

	import Searchbar from '$lib/components/admin/Searchbar.svelte';
	import MoreVerticalIcon from '@lucide/svelte/icons/more-vertical';

	export let data: PageData;

	const items = data.items;
</script>

<div class="flex flex-col gap-6">
	<h1 class="text-4xl font-medium">
		{data.collectionSlug.charAt(0).toUpperCase() + data.collectionSlug.slice(1)}
	</h1>


	<Searchbar id="search" name="search" class="max-w-md" />

	<div class="w-full overflow-hidden rounded-md border border-gray-200 bg-white">
		<table class="w-full divide-y divide-gray-200">
			<thead class="bg-gray-50">
				<tr>
					<th
						scope="col"
						class="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500"
					>
						Title
					</th>
					<th
						scope="col"
						class="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500"
					>
						Slug
					</th>
					<th
						scope="col"
						class="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500"
					>
						Author
					</th>
					<th
						scope="col"
						class="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500"
					>
						Updated
					</th>
					<th scope="col" class="relative px-6 py-3">
						<span class="sr-only">Actions</span>
					</th>
				</tr>
			</thead>
			<tbody class="divide-y divide-gray-200">
				{#each items as item}
					<tr class="hover:bg-gray-50">
						<td class="px-6 py-3">
							<a
								class="flex items-center text-sm text-gray-500"
								href={`/collections/${data.collectionSlug}/${item.resourceSlug}`}
							>
								{item.title}
							</a>
						</td>

						<td class="whitespace-nowrap px-6 py-3">
							<div class="flex items-center text-sm text-gray-500">
								{item.resourceSlug}
							</div>
						</td>

						<td class="whitespace-nowrap px-6 py-3">
							<div class="flex items-center text-sm text-gray-500">
								{item.created_by}
							</div>
						</td>

						<td class="whitespace-nowrap px-6 py-3">
							<div class="flex items-center text-sm text-gray-500">
								{new Date(item.updated_at).toLocaleString(undefined, {
									year: 'numeric',
									month: 'long',
									day: 'numeric',
									hour: '2-digit',
									minute: '2-digit'
								})}
							</div>
						</td>

						<td class="whitespace-nowrap px-6 py-3 text-right text-sm font-medium">
							<div class="flex items-center justify-end space-x-2">
								<button class="text-gray-400 hover:text-gray-600">
									<MoreVerticalIcon class="h-4 w-4" />
								</button>
							</div>
						</td>
					</tr>
				{/each}
			</tbody>
		</table>
	</div>
</div>
