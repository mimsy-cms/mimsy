<script lang="ts">
	import type { PageData } from './$types';

	import Searchbar from '$lib/components/admin/Searchbar.svelte';
	import MoreVerticalIcon from '@lucide/svelte/icons/more-vertical';
	import PlusIcon from '@lucide/svelte/icons/plus';

	export let data: PageData;

	const users = data.users;

	function formatDate(dateString: string): string {
		const date = new Date(dateString);
		return new Intl.DateTimeFormat('en-US', {
			year: 'numeric',
			month: 'short',
			day: 'numeric',
			hour: '2-digit',
			minute: '2-digit',
			hour12: false
		}).format(date);
	}
</script>

<div class="flex flex-col gap-6">
	<h1 class="text-4xl font-medium">Users</h1>

	<div class="flex justify-between gap-2">
		<Searchbar id="search" name="search" class="max-w-md" />

		<a
			href="/users/new"
			class="flex items-center rounded-md border border-gray-300 bg-blue-700 px-3 py-2 text-sm text-white hover:bg-blue-600"
		>
			<PlusIcon class="mr-3 h-5 w-5 flex-shrink-0" />
			Create
		</a>
	</div>

	<div class="w-full overflow-hidden rounded-md border border-gray-200 bg-white">
		<table class="w-full table-fixed divide-y divide-gray-200">
			<thead class="bg-gray-50">
				<tr>
					<th
						class="w-1/3 px-6 py-3 text-left text-xs font-medium tracking-wider text-gray-500 uppercase"
					>
						ID
					</th>
					<th
						class="w-1/3 px-6 py-3 text-left text-xs font-medium tracking-wider text-gray-500 uppercase"
					>
						Email
					</th>
					<th
						class="w-1/3 px-6 py-3 text-left text-xs font-medium tracking-wider text-gray-500 uppercase"
					>
						Updated
					</th>
					<th
						class="w-12 px-6 py-3 text-right text-xs font-medium tracking-wider text-gray-500 uppercase"
					>
						<span class="sr-only">Actions</span>
					</th>
				</tr>
			</thead>
			<tbody class="divide-y divide-gray-200">
				{#each users as user}
					<tr class="hover:bg-gray-50">
						<td class="w-1/3 px-6 py-3 text-sm text-gray-500">{user.id}</td>
						<td class="w-1/3 px-6 py-3 text-sm text-gray-500">{user.email}</td>
						<td class="w-1/3 px-6 py-3 text-sm text-gray-500">{formatDate(user.updated_at)}</td>
						<td class="w-12 px-6 py-3 text-right text-sm font-medium whitespace-nowrap">
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
