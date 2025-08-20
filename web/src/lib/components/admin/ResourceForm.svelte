<script lang="ts">
	import type {
		CollectionResource,
		Field,
		CollectionDefinition
	} from '$lib/collection/definition.js';
	import CheckboxField from '$lib/components/admin/fields/CheckboxField.svelte';
	import DateField from '$lib/components/admin/fields/DateField.svelte';
	import EmailField from '$lib/components/admin/fields/EmailField.svelte';
	import NumberField from '$lib/components/admin/fields/NumberField.svelte';
	import PlainTextField from '$lib/components/admin/fields/PlainTextField.svelte';
	import RichTextField from '$lib/components/admin/fields/RichTextField/RichTextField.svelte';
	import Input from '$lib/components/Input.svelte';
	import { onMount } from 'svelte';

	const {
		data,
		slugEditable = true
	}: {
		data: {
			definition: CollectionDefinition;
			resource?: CollectionResource;
		};
		slugEditable: boolean;
	} = $props();

	// Parse existing content and schema-defined fields
	let resourceContent: CollectionResource = $state({
		id: data.resource?.id || '',
		slug: data.resource?.slug || '',
		created_at: data.resource?.created_at || '',
		updated_at: data.resource?.updated_at || '',
		created_by: data.resource?.created_by || '',
		updated_by: data.resource?.updated_by || '',
		...Object.fromEntries(
			Object.keys(data.definition.fields).map((fieldName) => {
				const field = data.definition.fields[fieldName as string];
				let value = data.resource?.[fieldName] ?? getDefaultValue(field);
				if (field.type === 'date' && typeof value === 'string' && value) {
					value = new Date(value);
				} else if (field.type === 'rich_text') {
					if (!value) {
						value = null;
					}
				}
				return [fieldName, value];
			})
		)
	});

	let isSaving = $state(false);
	let error = $state('');
	let success = $state('');

	// eslint-disable-next-line @typescript-eslint/no-explicit-any
	function getDefaultValue(field: any) {
		switch (field.type) {
			case 'checkbox':
				return false;
			case 'number':
				return 0;
			case 'date':
				return new Date();
			case 'rich_text':
				return null;
			default:
				return '';
		}
	}

	function formatDate(dateStr: string): string {
		return new Date(dateStr).toLocaleString();
	}

	async function saveResource() {
		if (isSaving) return;

		try {
			isSaving = true;
			error = '';

			// eslint-disable-next-line @typescript-eslint/no-unused-vars
			const { id, created_at, updated_at, updated_by, slug, ...schemaContent }: CollectionResource =
				resourceContent;

			if (currentUser) {
				schemaContent.updated_by = currentUser.id;
			}

			Object.keys(data.definition.fields).forEach((fieldName) => {
				const field = data.definition.fields[fieldName as string];
				if (field.type === 'rich_text' && schemaContent[fieldName] !== undefined) {
					// Rich text should already be in the correct format from the editor
					// No need to transform it here since the editor handles the JSON structure
				}
			});

			const validationErrors: string[] = [];

			for (const [fieldName, field] of Object.entries(data.definition.fields)) {
				const value = schemaContent[fieldName];
				switch (field.type) {
					case 'number':
						if (typeof value !== 'number' || isNaN(value)) {
							validationErrors.push(`"${fieldName}" must be a number.`);
						}
						break;
					case 'date':
						if (!(value instanceof Date) || isNaN(value.getTime())) {
							validationErrors.push(`"${fieldName}" must be a date.`);
						}
						break;
					case 'rich_text':
						if (typeof value !== 'object' || value === null) {
							validationErrors.push(`"${fieldName}" must be text.`);
						}
						break;
					case 'checkbox':
						if (typeof value !== 'boolean') {
							validationErrors.push(`"${fieldName}" must be a boolean.`);
						}
						break;
					case 'email':
						if (typeof value !== 'string' || !/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(value)) {
							validationErrors.push(`"${fieldName}" must be an email address.`);
						}
						break;
					case 'string':
						if (typeof value !== 'string') {
							validationErrors.push(`"${fieldName}" must be text.`);
						}
				}
			}

			if (validationErrors.length > 0) {
				error = validationErrors.join('\n');
				isSaving = false;
				return;
			}

			const response = await fetch(
				`/api/v1/collections/${data.definition.slug}/${resourceContent.slug}`,
				{
					method: 'PUT',
					headers: {
						'Content-Type': 'application/json'
					},
					body: JSON.stringify(schemaContent)
				}
			);

			if (!response.ok) {
				throw new Error(`Failed to save resource: ${response.statusText}`) as Error;
			}

			const updatedResource = await response.json();

			resourceContent = {
				...resourceContent,
				...updatedResource
			};

			success = 'Resource saved successfully!';
		} catch (err) {
			console.error('Save error:', err);
			error = err instanceof Error ? err.message : 'Failed to save resource.';
		} finally {
			isSaving = false;
		}
	}

	function isRequired(field: Field): boolean {
		return isRequired(field) ?? false;
	}

	let currentUser: { id: string; email: string } | null = null;

	onMount(async () => {
		try {
			const res = await fetch('/api/v1/auth/me');
			if (!res.ok) {
				throw new Error('Not logged in');
			}

			currentUser = await res.json();
		} catch (error) {
			console.error('Error fetching current user:', error);
		}

		const response = await fetch('/api/v1/auth/me');
		if (response.ok) {
			currentUser = await response.json();
		}
	});
</script>

<div class="flex flex-col gap-6">
	<div class="flex items-center justify-between">
		<h1 class="text-4xl font-medium">{data.definition.name}</h1>
		<div class="flex gap-2">
			<button
				onclick={saveResource}
				class="rounded bg-blue-500 px-4 py-2 text-white hover:bg-blue-600 disabled:opacity-50"
				disabled={isSaving}
			>
				{isSaving ? 'Saving...' : 'Save Resource'}
			</button>
		</div>
	</div>

	{#if error}
		<div class="rounded border border-red-400 bg-red-100 px-4 py-3 text-red-700">
			<ul class="list-disc space-y-1 pl-5">
				{#each error.split('\n') as err (err)}
					<li>{err}</li>
				{/each}
			</ul>
		</div>
	{/if}

	{#if success}
		<div class="rounded border border-green-400 bg-green-100 px-4 py-3 text-green-700">
			{success}
		</div>
	{/if}

	<div class="flex gap-4">
		<div class="flex flex-1 flex-col gap-4 rounded-md border border-gray-300 bg-white p-4">
			<div class="flex flex-col gap-2">
				<label for="slug">Slug</label>
				<Input id="slug" name="slug" bind:value={resourceContent.slug} disabled={!slugEditable} />
			</div>

			{#if Object.keys(data.definition.fields).length > 0}
				{#each Object.entries(data.definition.fields) as [fieldName, field] (fieldName)}
					<div class="flex flex-col gap-2">
						{#if field.type === 'email'}
							<label for={fieldName}>
								{fieldName}
								{#if isRequired(field)}
									<span class="text-red-500">*</span>
								{/if}
							</label>
							<EmailField
								id={fieldName}
								name={fieldName}
								placeholder="example@example.com"
								bind:value={resourceContent[fieldName] as string}
							/>
						{:else if field.type === 'date'}
							<label for={fieldName}>
								{fieldName}
								{#if isRequired(field)}<span class="text-red-500">*</span>{/if}
							</label>
							<DateField
								id={fieldName}
								name={fieldName}
								label={field.label ?? field.name}
								bind:value={resourceContent[fieldName] as Date}
							/>
						{:else if field.type === 'number'}
							<label for={fieldName}>
								{fieldName}
								{#if isRequired(field)}<span class="text-red-500">*</span>{/if}
							</label>
							<NumberField
								id={fieldName}
								name={fieldName}
								bind:value={resourceContent[fieldName] as number}
							/>
						{:else if field.type === 'checkbox'}
							<label for={fieldName}>
								{fieldName}
								{#if isRequired(field)}<span class="text-red-500">*</span>{/if}
							</label>
							<CheckboxField
								id={fieldName}
								name={fieldName}
								label={field.label ?? field.name}
								bind:checked={resourceContent[fieldName] as boolean}
							/>
						{:else if field.type === 'rich_text'}
							<label for={fieldName}>
								{fieldName}
								{#if isRequired(field)}<span class="text-red-500">*</span>{/if}
							</label>
							<RichTextField bind:value={resourceContent[fieldName] as string | undefined} />
						{:else if field.type === 'string'}
							<label for={fieldName}>
								{fieldName}
								{#if isRequired(field)}<span class="text-red-500">*</span>{/if}
							</label>
							<PlainTextField
								id={fieldName}
								name={fieldName}
								bind:value={resourceContent[fieldName] as string}
							/>
						{:else}
							<p class="text-red-500">Unsupported field type: {field.type}</p>
						{/if}
					</div>
				{/each}
			{/if}
		</div>

		<div class="w-80 shrink-0 rounded-md border border-gray-300 bg-white p-4">
			<h2 class="text-2xl font-medium">Details</h2>
			<hr class="my-4 border-t-gray-300" />
			<div class="space-y-3 text-sm text-gray-700">
				<div>
					<p class="font-semibold">Slug</p>
					<p class="text-gray-600">/{resourceContent.slug}</p>
				</div>
				<div>
					<p class="font-semibold">Created</p>
					<p class="text-gray-600">
						{formatDate(data.resource?.created_at ?? new Date().toISOString())}
					</p>
				</div>
				<div>
					<p class="font-semibold">Created by</p>
					<p class="text-gray-600">{data.resource?.created_by_email ?? ''}</p>
				</div>
				<div>
					<p class="font-semibold">Last modified</p>
					<p class="text-gray-600">
						{formatDate(data.resource?.updated_at ?? new Date().toISOString())}
					</p>
				</div>
				<div>
					<p class="font-semibold">Last modified by</p>
					<p class="text-gray-600">{data.resource?.updated_by_email ?? ''}</p>
				</div>
			</div>
		</div>
	</div>
</div>
