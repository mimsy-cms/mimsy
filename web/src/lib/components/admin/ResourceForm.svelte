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
	import type { User } from '$lib/types/user';
	import type { SuperFormData } from 'sveltekit-superforms/client';
	import RelationshipField from './fields/RelationshipField.svelte';

	const {
		definition,
		resource,
		createdBy,
		updatedBy,
		slugEditable = true,
		form
	}: {
		definition: CollectionDefinition;
		resource?: CollectionResource;
		createdBy?: User;
		updatedBy?: User;
		slugEditable: boolean;
		form: SuperFormData<{
			[field: string]: any;
		}>;
	} = $props();

	function formatDate(dateStr: string): string {
		return new Date(dateStr).toLocaleString();
	}

	function isRequired(field: Field): boolean {
		return !!field.options.constraints.required;
	}
</script>

<div class="flex flex-col gap-6">
	<div class="flex items-center justify-between">
		<h1 class="text-4xl font-medium">{definition.name}</h1>
		<div class="flex gap-2">
			<button
				class="rounded bg-blue-500 px-4 py-2 text-white hover:bg-blue-600 disabled:opacity-50"
			>
				Save Resource
			</button>
		</div>
	</div>

	<div class="flex gap-4">
		<div class="flex flex-1 flex-col gap-4 rounded-md border border-gray-300 bg-white p-4">
			<div class="flex flex-col gap-2">
				<label for="slug">Slug</label>
				<Input id="slug" name="slug" bind:value={$form.slug} disabled={!slugEditable} />
			</div>

			{#if Object.keys(definition.fields).length > 0}
				{#each Object.entries(definition.fields) as [fieldName, field] (fieldName)}
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
								bind:value={$form[fieldName]}
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
								bind:value={$form[fieldName]}
							/>
						{:else if field.type === 'number'}
							<label for={fieldName}>
								{fieldName}
								{#if isRequired(field)}<span class="text-red-500">*</span>{/if}
							</label>
							<NumberField id={fieldName} name={fieldName} bind:value={$form[fieldName]} />
						{:else if field.type === 'checkbox'}
							<label for={fieldName}>
								{fieldName}
								{#if isRequired(field)}<span class="text-red-500">*</span>{/if}
							</label>
							<CheckboxField
								id={fieldName}
								name={fieldName}
								label={field.label ?? field.name}
								bind:checked={$form[fieldName]}
							/>
						{:else if field.type === 'rich_text'}
							<label for={fieldName}>
								{fieldName}
								{#if isRequired(field)}<span class="text-red-500">*</span>{/if}
							</label>
							<RichTextField bind:value={$form[fieldName]} />
						{:else if field.type === 'string'}
							<label for={fieldName}>
								{fieldName}
								{#if isRequired(field)}<span class="text-red-500">*</span>{/if}
							</label>
							<PlainTextField id={fieldName} name={fieldName} bind:value={$form[fieldName]} />
						{:else if field.type === 'relation'}
							{@const name = `${fieldName}_id`}
							<RelationshipField label={fieldName} {name} {field} bind:value={$form[fieldName]} />
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
					<p class="text-gray-600">/{$form.slug}</p>
				</div>
				<div>
					<p class="font-semibold">Created</p>
					<p class="text-gray-600">
						{formatDate(resource?.created_at ?? new Date().toISOString())}
					</p>
				</div>
				<div>
					<p class="font-semibold">Created by</p>
					<p class="text-gray-600">{createdBy?.email ?? ''}</p>
				</div>
				<div>
					<p class="font-semibold">Last modified</p>
					<p class="text-gray-600">
						{formatDate(resource?.updated_at ?? new Date().toISOString())}
					</p>
				</div>
				<div>
					<p class="font-semibold">Last modified by</p>
					<p class="text-gray-600">{updatedBy?.email ?? ''}</p>
				</div>
			</div>
		</div>
	</div>
</div>
