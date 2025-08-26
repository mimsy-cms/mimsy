<script lang="ts">
	import type { FieldRelation } from '$lib/collection/definition';
	import SelectField from './SelectField.svelte';
	
	type Props = {
		name: string;
		value?: string;
		label: string;
		field: FieldRelation;
	};
	
	type Item = {
		value: string;
		label: string;
		url?: string;
		contentType?: string;
	};
	
	function getEndpoint(resourceSlug: string) {
		switch (resourceSlug) {
			case '<builtins.user>':
				return '/api/v1/users';
			case '<builtins.media>':
				return '/api/v1/media';
			default:
				return `/api/v1/collections/${resourceSlug}`;
		}
	}
	
	function isImage(contentType: string): boolean {
		return contentType.startsWith('image/');
	}
	
	let { name, label, field, value = $bindable() }: Props = $props();
	let items = $state<Item[]>([]);
	let isOpen = $state(false);
	let selectedItem = $state<Item | null>(null);
	
	async function fetchResources() {
		const response = await fetch(getEndpoint(field.relatesTo));
		const resources = await response.json();
		
		switch (field.relatesTo) {
			case '<builtins.user>':
				// eslint-disable-next-line @typescript-eslint/no-explicit-any
				items = (resources ?? []).map((r: any) => ({ value: r.id, label: r.email }));
				break;
			case '<builtins.media>':
				// eslint-disable-next-line @typescript-eslint/no-explicit-any
				items = (resources ?? []).map((r: any) => ({ 
					value: r.id, 
					label: r.name,
					url: r.url,
					contentType: r.content_type
				}));
				break;
			default:
				// eslint-disable-next-line @typescript-eslint/no-explicit-any
				items = resources.map((r: any) => ({ value: r.id, label: r.slug }));
		}
		
		if (value) {
			selectedItem = items.find(item => item.value === value) || null;
		}
	}
	
	function handleSelect(item: Item) {
		value = item.value;
		selectedItem = item;
		isOpen = false;
	}
	
	function toggleDropdown() {
		isOpen = !isOpen;
	}
	
	$effect(() => {
		fetchResources();
	});
	
	$effect(() => {
		if (value && items.length > 0) {
			selectedItem = items.find(item => item.value === value) || null;
		}
	});
</script>

{#if field.relatesTo === '<builtins.media>'}
	<div class="flex flex-col gap-2">
		<label for={name} class="block text-sm font-medium text-gray-700">
			{label}
		</label>
		
		<div class="relative">
			<button
				type="button"
				class="relative w-full cursor-default rounded-md border border-gray-300 bg-white py-2 pl-3 pr-10 text-left shadow-sm focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500 sm:text-sm"
				onclick={toggleDropdown}
			>
				<span class="flex items-center">
					{#if selectedItem}
						{#if selectedItem.url && selectedItem.contentType && isImage(selectedItem.contentType)}
							<img
								src={selectedItem.url}
								alt={selectedItem.label}
								class="h-6 w-6 flex-shrink-0 rounded object-cover"
							/>
						{:else}
							<div class="h-6 w-6 flex-shrink-0 rounded bg-gray-200 flex items-center justify-center">
								<svg class="h-4 w-4 text-gray-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
								</svg>
							</div>
						{/if}
						<span class="ml-3 block truncate">{selectedItem.label}</span>
					{:else}
						<span class="block truncate text-gray-500">Select media...</span>
					{/if}
				</span>
				<span class="pointer-events-none absolute inset-y-0 right-0 ml-3 flex items-center pr-2">
					<svg class="h-5 w-5 text-gray-400" viewBox="0 0 20 20" fill="currentColor">
						<path fill-rule="evenodd" d="M10 3a1 1 0 01.707.293l3 3a1 1 0 01-1.414 1.414L10 5.414 7.707 7.707a1 1 0 01-1.414-1.414l3-3A1 1 0 0110 3z" clip-rule="evenodd" transform="rotate(180 10 10)" />
					</svg>
				</span>
			</button>
			
			{#if isOpen}
				<div class="absolute z-10 mt-1 w-full rounded-md bg-white shadow-lg ring-1 ring-black ring-opacity-5 max-h-60 overflow-auto">
					<ul class="py-1">
						{#each items as item}
							<button
								type="button"
								class="relative w-full text-left select-none py-2 pl-3 pr-9 hover:bg-blue-50 cursor-pointer"
								onclick={() => handleSelect(item)}
								onkeydown={(e) => e.key === 'Enter' && handleSelect(item)}
							>
								<div class="flex items-center">
									{#if item.url && item.contentType && isImage(item.contentType)}
										<img
											src={item.url}
											alt={item.label}
											class="h-8 w-8 flex-shrink-0 rounded object-cover"
										/>
									{:else}
										<div class="h-8 w-8 flex-shrink-0 rounded bg-gray-200 flex items-center justify-center">
											<svg class="h-5 w-5 text-gray-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
												<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
											</svg>
										</div>
									{/if}
									<div class="ml-3 flex-1 min-w-0">
										<span class="block truncate font-medium text-gray-900">
											{item.label}
										</span>
										{#if item.contentType}
											<span class="block truncate text-xs text-gray-500">
												{item.contentType}
											</span>
										{/if}
									</div>
								</div>
								
								{#if selectedItem?.value === item.value}
									<span class="absolute inset-y-0 right-0 flex items-center pr-4">
										<svg class="h-5 w-5 text-blue-600" viewBox="0 0 20 20" fill="currentColor">
											<path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd" />
										</svg>
									</span>
								{/if}
							</button>
						{/each}
					</ul>
				</div>
			{/if}
		</div>
		
		<input type="hidden" {name} {value} />
	</div>
{:else}
	<SelectField {name} {label} {items} bind:value />
{/if}

<!-- Click outside to close dropdown -->
{#if isOpen}
	<button
		type="button"
		class="fixed inset-0 z-0"
		onclick={() => isOpen = false}
		onkeydown={(e) => e.key === 'Escape' && (isOpen = false)}
		aria-label="Close dropdown"
	></button>
{/if}