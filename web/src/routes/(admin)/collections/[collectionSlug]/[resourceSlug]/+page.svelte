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

	let { data } = $props();

	// Parse existing content and schema-defined fields
	let resourceContent = $state({
		title: data.resource?.title || data.resource?.content?.title || '',
		slug: data.resource?.slug || '',
		...Object.fromEntries(
			Object.keys(data.definition.fields).map((fieldName) => [
				fieldName, 
				data.resource?.content?.[fieldName] || data.resource?.[fieldName] || getDefaultValue(data.definition.fields[fieldName])
			])
		)
	});

	// Track any additional fields that exist in content but not in schema
	let additionalFields = $state({});

	// Initialize additional fields from existing content
	$effect(() => {
		if (data.resource?.content) {
			const schemaFieldNames = new Set([
				'title', 'slug', 
				...Object.keys(data.definition.fields)
			]);
			
			const additionalContentFields = {};
			for (const [key, value] of Object.entries(data.resource.content)) {
				if (!schemaFieldNames.has(key)) {
					additionalContentFields[key] = {
						value: value,
						type: guessFieldType(value)
					};
				}
			}
			additionalFields = additionalContentFields;
		}
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

	function addContentField() {
		if (!newFieldName.trim()) {
			error = 'Field name is required';
			return;
		}

		const fieldKey = newFieldName.toLowerCase().replace(/\s+/g, '_');
		
		// Check if field already exists
		if (fieldKey in resourceContent || fieldKey in additionalFields) {
			error = 'Field with this name already exists';
			return;
		}

		// Convert value based on type
		let value = newFieldValue;
		if (newFieldType === 'number') {
			value = parseFloat(newFieldValue) || 0;
		} else if (newFieldType === 'checkbox') {
			value = newFieldValue === 'true' || newFieldValue === '1';
		}

		additionalFields[fieldKey] = {
			value: value,
			type: newFieldType
		};

		// Reset form
		newFieldName = '';
		newFieldValue = '';
		newFieldType = 'plaintext';
		showAddFieldForm = false;
		error = '';
		success = 'Field added successfully';
		clearMessages();
	}

	function removeContentField(fieldName: string) {
		if (confirm(`Are you sure you want to remove the "${fieldName}" field?`)) {
			delete additionalFields[fieldName];
			additionalFields = { ...additionalFields };
			success = 'Field removed successfully';
			clearMessages();
		}
	}

	function updateAdditionalField(fieldName: string, value: any) {
		additionalFields[fieldName].value = value;
	}

	async function saveResource() {
		if (isSaving) return;

		try {
			isSaving = true;
			error = '';

			// Combine all content data
			const { id, created_at, updated_at, ...schemaContent } = resourceContent;
			const additionalContent = Object.fromEntries(
				Object.entries(additionalFields).map(([key, field]) => [key, field.value])
			);
			
			const contentData = { ...schemaContent, ...additionalContent };

			const response = await fetch(`/api/v1/collections/${data.definition.slug}/${data.resource.slug}`, {
				method: 'PUT',
				headers: {
					'Content-Type': 'application/json'
				},
				body: JSON.stringify(contentData)
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
				<label for="title">Title</label>
				<Input
					type="text"
					id="title"
					name="title"
					bind:value={resourceContent.title}
				/>
			</div>

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
							<RichTextField 
								bind:value={resourceContent[fieldName]} 
							/>
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

			<!-- Additional Content Fields -->
			<div class="pt-4">
				<div class="flex items-center justify-between mb-4">
					<h3 class="text-lg font-medium">Content</h3>
					<button
						onclick={() => showAddFieldForm = !showAddFieldForm}
						class="rounded bg-green-500 px-3 py-1 text-white hover:bg-green-600 text-sm"
					>
						{showAddFieldForm ? 'Cancel' : 'Add Field'}
					</button>
				</div>

				<!-- Add Field Form -->
				{#if showAddFieldForm}
					<div class="rounded-lg bg-gray-50 p-4 mb-4">
						<h4 class="text-md font-medium mb-3">Add New Content Field</h4>
						<div class="grid grid-cols-1 md:grid-cols-3 gap-3">
							<div>
								<label for="newFieldName" class="block text-sm font-medium mb-1">Field Name</label>
								<Input
									id="newFieldName"
									bind:value={newFieldName}
									placeholder="field_name"
								/>
							</div>
							<div>
								<label for="newFieldType" class="block text-sm font-medium mb-1">Field Type</label>
								<select 
									id="newFieldType" 
									bind:value={newFieldType}
									class="w-full rounded border border-gray-300 px-3 py-2"
								>
									{#each fieldTypes as type}
										<option value={type.value}>{type.label}</option>
									{/each}
								</select>
							</div>
							<div>
								<label for="newFieldValue" class="block text-sm font-medium mb-1">
									{newFieldType === 'checkbox' ? 'Value (true/false)' : 'Initial Value'}
								</label>
								{#if newFieldType === 'checkbox'}
									<select 
										bind:value={newFieldValue}
										class="w-full rounded border border-gray-300 px-3 py-2"
									>
										<option value="false">False</option>
										<option value="true">True</option>
									</select>
								{:else}
									<Input
										id="newFieldValue"
										bind:value={newFieldValue}
										placeholder="Enter value"
										type={newFieldType === 'number' ? 'number' : newFieldType === 'date' ? 'date' : 'text'}
									/>
								{/if}
							</div>
						</div>
						<div class="mt-3">
							<button
								onclick={addContentField}
								class="rounded bg-blue-500 px-4 py-2 text-white hover:bg-blue-600 text-sm"
							>
								Add Field
							</button>
						</div>
					</div>
				{/if}

				<!-- Existing Additional Fields -->
				{#if Object.keys(additionalFields).length === 0}
					<p class="text-gray-500 text-sm">No additional content fields. Click "Add Field" to create one.</p>
				{:else}
					{#each Object.entries(additionalFields) as [fieldName, field] (fieldName)}
						<div class="flex flex-col gap-2 mb-4">
							<div class="flex items-center justify-between">
								<label for={fieldName} class="font-medium">
									{fieldName.replace(/_/g, ' ').replace(/\b\w/g, l => l.toUpperCase())}
									<span class="text-xs text-gray-500 font-normal">({field.type})</span>
								</label>
								<button
									onclick={() => removeContentField(fieldName)}
									class="rounded bg-red-500 px-2 py-1 text-white hover:bg-red-600 text-xs"
								>
									Remove
								</button>
							</div>
							
							{#if field.type === 'checkbox'}
								<label class="flex items-center gap-2">
									<input 
										type="checkbox" 
										bind:checked={field.value}
										onchange={() => updateAdditionalField(fieldName, field.value)}
									/>
									<span class="text-sm">{fieldName.replace(/_/g, ' ')}</span>
								</label>
							{:else if field.type === 'number'}
								<Input
									type="number"
									bind:value={field.value}
									onchange={() => updateAdditionalField(fieldName, field.value)}
								/>
							{:else if field.type === 'date'}
								<Input
									type="date"
									bind:value={field.value}
									onchange={() => updateAdditionalField(fieldName, field.value)}
								/>
							{:else}
								<Input
									type="text"
									bind:value={field.value}
									onchange={() => updateAdditionalField(fieldName, field.value)}
								/>
							{/if}
						</div>
					{/each}
				{/if}
			</div>
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
					<p class="text-gray-600">{formatDate(data.definition.updated_at)}</p>
				</div>
				<div>
					<p class="font-semibold">Last modified by</p>
					<p class="text-gray-600">{data.definition.updated_by || 'N/A'}</p>
				</div>
			</div>

			<!-- Content Summary -->
			<div class="mt-6">
				<h3 class="text-base font-medium">Content Fields</h3>
				<hr class="my-2 border-t-gray-300" />
				<div class="space-y-1 text-xs">
					<div class="text-gray-600">
						<strong>Schema Fields:</strong> {Object.keys(data.definition.fields).length}
					</div>
					<div class="text-gray-600">
						<strong>Additional Fields:</strong> {Object.keys(additionalFields).length}
					</div>
					<div class="text-gray-600">
						<strong>Total Fields:</strong> {Object.keys(data.definition.fields).length + Object.keys(additionalFields).length}
					</div>
				</div>
			</div>
		</div>
	</div>
</div>