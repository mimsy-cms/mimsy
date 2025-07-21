<script lang="ts">
	import CheckboxField from '$lib/components/admin/fields/CheckboxField.svelte';
	import DateField from '$lib/components/admin/fields/DateField.svelte';
	import EmailField from '$lib/components/admin/fields/EmailField.svelte';
	import NumberField from '$lib/components/admin/fields/NumberField.svelte';
	import PlainTextField from '$lib/components/admin/fields/PlainTextField.svelte';
	import RelationshipField from '$lib/components/admin/fields/RelationshipField.svelte';
	import RichTextField from '$lib/components/admin/fields/RichTextField.svelte';
	import SelectField from '$lib/components/admin/fields/SelectField.svelte';

	let plainTextValue = $state('');
</script>

<div class="flex flex-col gap-6">
	<h1 class="text-4xl font-medium">Pages</h1>

	<div class="flex gap-4">
		<div class="flex-1 rounded-md border border-gray-300 bg-white p-4">
			<div class="flex flex-col gap-2">
				<label for="plainTextField">PlainTextField</label>
				<PlainTextField
					id="plainTextField"
					name="plainTextField"
					bind:value={plainTextValue}
					placeholder="This is the placeholder"
				/>
			</div>

			<RichTextField />
			<EmailField />
			<NumberField />
			<DateField />
			<CheckboxField />
			<SelectField />
			<RelationshipField />
		</div>

		<div class="rounded-md border border-gray-300 bg-white p-4">
			<h2 class="text-2xl font-medium">Details</h2>
			<hr class="my-4 border-t-gray-300" />
		</div>
	</div>
</div>
