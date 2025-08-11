<script lang="ts">
	import CheckboxField from '$lib/components/admin/fields/CheckboxField.svelte';
	import DateField from '$lib/components/admin/fields/DateField.svelte';
	import EmailField from '$lib/components/admin/fields/EmailField.svelte';
	import NumberField from '$lib/components/admin/fields/NumberField.svelte';
	import PlainTextField from '$lib/components/admin/fields/PlainTextField.svelte';
	import RichTextField from '$lib/components/admin/fields/RichTextField/RichTextField.svelte';
	import SelectField from '$lib/components/admin/fields/SelectField.svelte';

	let { data } = $props();

	let fields = $derived(Object.entries(data.definition.fields));

	function formatDate(dateStr: string): string {
		return new Date(dateStr).toLocaleString();
	}
</script>

<div class="flex flex-col gap-6">
	<h1 class="text-4xl font-medium">{data.definition.name}</h1>

	<div class="flex gap-4">
		<div class="flex flex-1 flex-col gap-4 rounded-md border border-gray-300 bg-white p-4">
			{#each fields as [fieldName, field] (fieldName)}
				<div class="flex flex-col gap-2">
					{#if field.type === 'email'}
						<label for={field.name}>
							{field.label}
						</label>
						<EmailField id={field.name} name={field.name} placeholder="example@example.com" />
					{:else if field.type === 'date'}
						<label for={field.name}>
							{field.label}
						</label>
						<DateField id={field.name} name={field.name} label={field.label} />
					{:else if field.type === 'number'}
						<label for={field.name}>
							{field.label}
						</label>
						<NumberField id={field.name} name={field.name} />
					{:else if field.type === 'checkbox'}
						<CheckboxField id={field.name} name={field.name} label={field.label} />
					{:else if field.type === 'select'}
						<label for={field.name}>
							{field.label}
						</label>
						<SelectField name={field.name} options={field.options} />
					{:else if field.type === 'richtext'}
						<label for={field.name}>
							{field.label}
						</label>
						<RichTextField />
					{:else if field.type === 'plaintext'}
						<label for={field.name}>
							{field.label}
						</label>
						<PlainTextField id={field.name} name={field.name} />
					{:else}
						<p class="text-red-500">Unsupported field type: {field.type}</p>
					{/if}
				</div>
			{/each}
		</div>

		<div class="w-80 shrink-0 rounded-md border border-gray-300 bg-white p-4">
			<h2 class="text-2xl font-medium">Details</h2>
			<hr class="my-4 border-t-gray-300" />

			<div class="space-y-1 text-sm text-gray-700">
				<p class="text-base"><strong>Slug</strong></p>
				<p>/{data.resource.slug}</p>

				<p class="text-base"><strong>Created</strong></p>
				<p>{formatDate(data.definition.created_at)}</p>

				<p class="text-base"><strong>Created by</strong></p>
				<p>{data.definition.created_by}</p>

				<p class="text-base"><strong>Last modified</strong></p>
				<p>{formatDate(data.definition.updated_at)}</p>

				<p class="text-base"><strong>Last modified by</strong></p>
				<p>{data.definition.updated_by}</p>
			</div>
		</div>
	</div>
</div>
