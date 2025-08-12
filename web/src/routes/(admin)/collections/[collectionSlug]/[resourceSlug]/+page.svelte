<script lang="ts">
	import CheckboxField from '$lib/components/admin/fields/CheckboxField.svelte';
	import DateField from '$lib/components/admin/fields/DateField.svelte';
	import EmailField from '$lib/components/admin/fields/EmailField.svelte';
	import NumberField from '$lib/components/admin/fields/NumberField.svelte';
	import PlainTextField from '$lib/components/admin/fields/PlainTextField.svelte';
	import RichTextField from '$lib/components/admin/fields/RichTextField/RichTextField.svelte';
	import SelectField from '$lib/components/admin/fields/SelectField.svelte';
	import Input from '$lib/components/Input.svelte';
	import { env } from '$env/dynamic/public';
	import { fi } from 'zod/v4/locales';

	let { data } = $props();

	// Parse existing content and schema-defined fields
	let resourceContent = $state({
		slug: data.resource?.slug || '',
		...Object.fromEntries(
			Object.keys(data.definition.fields).map((fieldName) => [
				fieldName, data.resource?.content?.[fieldName] ?? getDefaultValue(data.definition.fields[fieldName])
			])
		)
	});

	// Form for adding new content fields
	let showAddFieldForm = $state(false);
	let newFieldName = $state('');
	let newFieldValue = $state('');
	let newFieldType = $state('plaintext');

	let isSaving = $state(false);
	let error = $state('');
	let success = $state('');

	const fieldTypes = [
		{ value: 'plaintext', label: 'Text' },
		{ value: 'number', label: 'Number' },
		{ value: 'checkbox', label: 'Checkbox' },
		{ value: 'date', label: 'Date' },
		{ value: 'email', label: 'Email' },
		{ value: 'richtext', label: 'Rich Text' },
		{ value: 'select', label: 'Select' }
	];

	function getDefaultValue(field: any) {
		switch (field.type) {
			case 'checkbox':
				return false;
			case 'number':
				return 0;
			case 'date':
				return '';
			default:
				return '';
		}
	}

	function guessFieldType(value: any): string {
		if (typeof value === 'boolean') return 'checkbox';
		if (typeof value === 'number') return 'number';
		if (typeof value === 'string' && value.match(/^\d{4}-\d{2}-\d{2}/)) return 'date';
		return 'plaintext';
	}

	function formatDate(dateStr: string): string {
		return new Date(dateStr).toLocaleString();
	}

	async function saveResource() {
		if (isSaving) return;

		try {
			isSaving = true;
			error = '';

			const { id, created_at, updated_at, slug, ...schemaContent } = resourceContent;

			const response = await fetch(`/api/v1/collections/${data.definition.slug}/${resourceContent.slug}`, {
				method: 'PUT',
				headers: {
					'Content-Type': 'application/json'
				},
				body: JSON.stringify(schemaContent)
			});

			if (!response.ok) {
				throw new Error(`Failed to save resource: ${response.statusText}`);
			}

			const updatedResource = await response.json();

			resourceContent = {
				...resourceContent,
				...updatedResource
			};

			success = 'Resource saved successfully!';
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to save resource.';
		} finally {
			isSaving = false;
		}
	}

	function clearMessages() {
		setTimeout(() => {
			error = '';
			success = '';
		}, 5000);
	}
</script>

<div class="flex flex-col gap-6">
	<!-- Header with actions -->
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

	<!-- Success/Error Messages -->
	{#if error}
		<div class="rounded bg-red-100 border border-red-400 text-red-700 px-4 py-3">
			{error}
		</div>
	{/if}

	{#if success}
		<div class="rounded bg-green-100 border border-green-400 text-green-700 px-4 py-3">
			{success}
		</div>
	{/if}

	<div class="flex gap-4">
		<div class="flex flex-1 flex-col gap-4 rounded-md border border-gray-300 bg-white p-4">
			<div class="flex flex-col gap-2">
				<label for="slug">Slug</label>
				<Input
					id="slug"
					name="slug"
					bind:value={resourceContent.slug}
				/>
			</div>

			{#if Object.keys(data.definition.fields).length > 0}
				<div class="flex flex-col gap-2">
					<label class="text-lg font-medium">Schema Fields</label>
				</div>

				{#each Object.entries(data.definition.fields) as [fieldName, field] (fieldName)}
					<div class="flex flex-col gap-2">
						{#if field.type === 'email'}
							<label for={fieldName}>
								{field.label}
								{#if field.required}<span class="text-red-500">*</span>{/if}
							</label>
							<EmailField 
								id={fieldName} 
								name={fieldName} 
								placeholder="example@example.com"
								bind:value={resourceContent[fieldName]}
								required={field.required}
							/>
						{:else if field.type === 'date'}
							<label for={fieldName}>
								{field.label}
								{#if field.required}<span class="text-red-500">*</span>{/if}
							</label>
							<DateField 
								id={fieldName} 
								name={fieldName} 
								label={field.label}
								bind:value={resourceContent[fieldName]}
								required={field.required}
							/>
						{:else if field.type === 'number'}
							<label for={fieldName}>
								{field.label}
								{#if field.required}<span class="text-red-500">*</span>{/if}
							</label>
							<NumberField 
								id={fieldName} 
								name={fieldName}
								bind:value={resourceContent[fieldName]}
								required={field.required}
							/>
						{:else if field.type === 'checkbox'}
							<CheckboxField 
								id={fieldName} 
								name={fieldName} 
								label={field.label}
								bind:checked={resourceContent[fieldName]}
								required={field.required}
							/>
						{:else if field.type === 'select'}
							<label for={fieldName}>
								{field.label}
								{#if field.required}<span class="text-red-500">*</span>{/if}
							</label>
							<SelectField 
								name={fieldName} 
								options={field.options || []}
								bind:value={resourceContent[fieldName]}
								required={field.required}
							/>
						{:else if field.type === 'richtext'}
							<label for={fieldName}>
								{field.label}
								{#if field.required}<span class="text-red-500">*</span>{/if}
							</label>
							<!-- <RichTextField
								bind:value={resourceContent[fieldName]} 
							/> -->
						{:else if field.type === 'plaintext'}
							<label for={fieldName}>
								{field.label}
								{#if field.required}<span class="text-red-500">*</span>{/if}
							</label>
							<PlainTextField 
								id={fieldName} 
								name={fieldName}
								bind:value={resourceContent[fieldName]}
								required={field.required}
							/>
						{:else}
							<p class="text-red-500">Unsupported field type: {field.type}</p>
						{/if}
					</div>
				{/each}
			{/if}
		</div>

		<!-- Details Sidebar -->
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
					<p class="text-gray-600">{formatDate(data.definition.created_at)}</p>
				</div>
				<div>
					<p class="font-semibold">Created by</p>
					<p class="text-gray-600">{data.definition.created_by}</p>
				</div>
				<div>
					<p class="font-semibold">Last modified</p>
					<p class="text-gray-600">{formatDate(data.resource.updated_at)}</p>
				</div>
				<div>
					<p class="font-semibold">Last modified by</p>
					<p class="text-gray-600">{data.definition.updated_by || 'N/A'}</p>
				</div>
			</div>
		</div>
	</div>
</div>