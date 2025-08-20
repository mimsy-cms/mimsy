<script lang="ts">
	import Searchbar from '$lib/components/admin/Searchbar.svelte';
	import Plus from '@lucide/svelte/icons/plus';
	import { enhance } from '$app/forms';

	let { data, form } = $props();

	let showCreateModal = $state(false);
	let newResourceSlug = $state('');
	let isCreating = $state(false);

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

	function openCreateModal() {
		showCreateModal = true;
		newResourceSlug = '';
	}

	function closeCreateModal() {
		showCreateModal = false;
		newResourceSlug = '';
	}

	function handleKeydown(event: KeyboardEvent) {
		if (event.key === 'Escape') {
			closeCreateModal();
		}
	}
	
	$effect(() => {
		if (form?.error) {
			// Keep modal open to show error
		} else if (form !== null && !form?.error) {
			showCreateModal = false;
			newResourceSlug = '';
		}
	});
</script>

<div class="flex flex-col gap-6">
	<h1 class="text-4xl font-medium">
		{data.collectionName}
	</h1>

	<div class="flex items-center justify-between">
		<Searchbar id="search" name="search" class="max-w-md" />
		
		<button
			onclick={openCreateModal}
			class="inline-flex items-center gap-2 px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 transition-colors"
		>
			<Plus size={16} />
			Create
		</button>
	</div>

	<div class="w-full overflow-hidden rounded-md border border-gray-200 bg-white">
		<table class="w-full table-fixed divide-y divide-gray-200">
			<thead class="bg-gray-50">
				<tr>
					<th class="w-1/4 px-6 py-3 text-left text-xs font-medium tracking-wider text-gray-500 uppercase">
						ID
					</th>
					<th class="w-1/4 px-6 py-3 text-left text-xs font-medium tracking-wider text-gray-500 uppercase">
						Slug
					</th>
					<th class="w-1/4 px-6 py-3 text-left text-xs font-medium tracking-wider text-gray-500 uppercase">
						Created By
					</th>
					<th class="w-1/4 px-6 py-3 text-left text-xs font-medium tracking-wider text-gray-500 uppercase">
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
				{:else}
					<tr>
						<td colspan="4" class="px-6 py-8 text-center text-gray-500">
							<div class="flex flex-col items-center gap-2">
								<p>No resources found</p>
								<button
									onclick={openCreateModal}
									class="text-blue-600 hover:text-blue-800 hover:underline"
								>
									Create your first resource
								</button>
							</div>
						</td>
					</tr>
				{/each}
			</tbody>
		</table>
	</div>
</div>

<!-- Create Resource Modal -->
{#if showCreateModal}
	<div class="fixed inset-0 z-50 overflow-y-auto" aria-labelledby="modal-title" role="dialog" aria-modal="true">
		<!-- Background overlay -->
		<div class="flex items-end justify-center min-h-screen pt-4 px-4 pb-20 text-center sm:block sm:p-0">
			<div 
				class="fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity" 
				aria-hidden="true"
				onclick={closeCreateModal}
			></div>

			<!-- This element is to trick the browser into centering the modal contents. -->
			<span class="hidden sm:inline-block sm:align-middle sm:h-screen" aria-hidden="true">&#8203;</span>

			<!-- Modal panel -->
			<div class="relative inline-block align-bottom bg-white rounded-lg px-4 pt-5 pb-4 text-left overflow-hidden shadow-xl transform transition-all sm:my-8 sm:align-middle sm:max-w-lg sm:w-full sm:p-6">
				<div>
					<div class="mx-auto flex items-center justify-center h-12 w-12 rounded-full bg-blue-100">
						<Plus class="h-6 w-6 text-blue-600" aria-hidden="true" />
					</div>
					<div class="mt-3 text-center sm:mt-5">
						<h3 class="text-lg leading-6 font-medium text-gray-900" id="modal-title">
							Create New Resource
						</h3>
						<div class="mt-2">
							<p class="text-sm text-gray-500">
								Enter a unique slug for the new resource in "{data.collectionName}".
							</p>
						</div>
					</div>
				</div>
				
				<form 
					method="POST" 
					action="?/create"
					use:enhance={({ formElement, formData, action, cancel }) => {
						isCreating = true;
						return async ({ result, update }) => {
							isCreating = false;
							// Let SvelteKit handle the result
							await update();
						};
					}}
				>
					<div class="mt-5 sm:mt-6">
						<div class="mb-4">
							<label for="resource-slug" class="block text-sm font-medium text-gray-700 mb-2">
								Resource Slug
							</label>
							<input
								id="resource-slug"
								name="slug"
								type="text"
								bind:value={newResourceSlug}
								onkeydown={handleKeydown}
								placeholder="e.g., my-new-resource"
								required
								pattern="^[a-z0-9]+(?:-[a-z0-9]+)*$"
								class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
								class:border-red-300={form?.error}
								class:focus:ring-red-500={form?.error}
								class:focus:border-red-500={form?.error}
							/>
							{#if form?.error}
								<p class="mt-1 text-sm text-red-600">{form.error}</p>
							{/if}
							<p class="mt-1 text-xs text-gray-500">
								Use lowercase letters, numbers, and hyphens only. No spaces or special characters.
							</p>
						</div>
						
						<div class="flex gap-3 sm:flex-row-reverse">
							<button
								type="submit"
								disabled={isCreating || !newResourceSlug.trim()}
								class="w-full inline-flex justify-center rounded-md border border-transparent shadow-sm px-4 py-2 bg-blue-600 text-base font-medium text-white hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 sm:ml-3 sm:w-auto sm:text-sm disabled:opacity-50 disabled:cursor-not-allowed"
							>
								{#if isCreating}
									<div class="flex items-center gap-2">
										<div class="animate-spin rounded-full h-4 w-4 border-b-2 border-white"></div>
										Creating...
									</div>
								{:else}
									Create Resource
								{/if}
							</button>
							
							<button
								type="button"
								onclick={closeCreateModal}
								disabled={isCreating}
								class="mt-3 w-full inline-flex justify-center rounded-md border border-gray-300 shadow-sm px-4 py-2 bg-white text-base font-medium text-gray-700 hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 sm:mt-0 sm:w-auto sm:text-sm disabled:opacity-50 disabled:cursor-not-allowed"
							>
								Cancel
							</button>
						</div>
					</div>
				</form>
			</div>
		</div>
	</div>
{/if}