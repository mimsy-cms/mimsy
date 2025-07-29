<script lang="ts">
	import { onMount } from 'svelte';
	import type { PageData } from './$types';

	import CheckboxField from '$lib/components/admin/fields/CheckboxField.svelte';
	import DateField from '$lib/components/admin/fields/DateField.svelte';
	import EmailField from '$lib/components/admin/fields/EmailField.svelte';
	import NumberField from '$lib/components/admin/fields/NumberField.svelte';
	import PlainTextField from '$lib/components/admin/fields/PlainTextField.svelte';
	import RichTextField from '$lib/components/admin/fields/RichTextField/RichTextField.svelte';
	import SelectField from '$lib/components/admin/fields/SelectField.svelte';
	import { email } from 'zod/v4';

	// let plainTextValue = $state('');
	// let emailValue = $state('');
	// let numberValue = $state<number>();
	// let dateValue = $state<Date>();
	// let selectValue = $state<string>();

	export let data: PageData;

	let values: Record<string, any> = {};

	onMount(() => {
		for (const field of data.definition.fields) {
			values[field.name] = null;
		}
	});

	const fieldComponents = {
		plain_text: PlainTextField,
		email: EmailField,
		number: NumberField,
		date: DateField,
		checkbox: CheckboxField,
		select: SelectField,
		rich_text: RichTextField
	};
</script>


<div class="flex flex-col gap-6">
	<h1 class="text-4xl font-medium">{data.definition.name}</h1>

	<div class="flex gap-4">
		<div class="flex flex-1 flex-col gap-4 rounded-md border border-gray-300 bg-white p-4">
			{#each data.definition.fields as field (field.name)}
				<div class="flex flex-col gap-2">
					{#if fieldComponents[field.type]}
						<label for={field.name}>{field.label}</label>
						<svelte:component
							this={fieldComponents[field.type]}
							id={field.name}
							name={field.name}
							bind:value={values[field.name]}
							placeholder={field.placeholder}
							label={field.label}
							items={field.items}
							multiple={field.multiple}
						/>
					{:else}
						<p class="text-red-500">Unsupported field type: {field.type}</p>
					{/if}
				</div>
			{/each}
		</div>

		<div class="rounded-md border border-gray-300 bg-white p-4">
			<h2 class="text-2xl font-medium">Details</h2>
			<hr class="my-4 border-t-gray-300" />
			<pre class="text-sm text-gray-500">{JSON.stringify(values, null, 2)}</pre>
		</div>
	</div>
</div>

<!-- <div class="flex flex-col gap-6">
	<h1 class="text-4xl font-medium">Pages</h1>

	<div class="flex gap-4">
		<div class="flex flex-1 flex-col gap-2 rounded-md border border-gray-300 bg-white p-4">
			<div class="flex flex-col gap-2">
				<label for="plainTextField">PlainTextField</label>
				<PlainTextField
					id="plainTextField"
					name="plainTextField"
					bind:value={plainTextValue}
					placeholder="This is the placeholder"
				/>
			</div>

			<div class="flex flex-col gap-2">
				<label for="richTextField">RichTextField</label>
				<RichTextField />
			</div>

			<div class="flex flex-col gap-2">
				<label for="emailField">EmailField</label>
				<EmailField id="emailField" name="emailField" bind:value={emailValue} />
			</div>

			<div class="flex flex-col gap-2">
				<label for="numberField">NumberField</label>
				<NumberField id="numberField" name="numberField" bind:value={numberValue} />
			</div>

			<div class="flex flex-col gap-2">
				<DateField id="dateField" name="dateField" bind:value={dateValue} label={'DateField'} />
			</div>

			<div class="flex flex-col gap-2">
				<CheckboxField id="checkboxField" name="checkboxField" label={'CheckboxField'} />
			</div>

			<div class="flex flex-col gap-2">
				<label for="selectField">SelectField</label>
				<SelectField
					name="selectField"
					bind:value={selectValue}
					label="Select something"
					items={[
						{ value: 'option1', label: 'Option 1' },
						{ value: 'option2', label: 'Option 2' },
						{ value: 'option3', label: 'Option 3' }
					]}
					multiple
				/>
			</div>
		</div>

		<div class="rounded-md border border-gray-300 bg-white p-4">
			<h2 class="text-2xl font-medium">Details</h2>
			<hr class="my-4 border-t-gray-300" />
		</div>
	</div>
</div> -->
