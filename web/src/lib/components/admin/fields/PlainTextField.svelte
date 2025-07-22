<script lang="ts">
	type Props = {
		id: string;
		name: string;
		value: string;
		placeholder?: string;
		disabled?: boolean;
		multiline?: boolean;
	};

	let { id, name, value = $bindable(), placeholder, disabled, multiline }: Props = $props();

	// We want the number of rows to be at least 3 and grow based on the value of the textarea.
	let rows = $derived(Math.max(value?.split('\n').length, 3));
</script>

{#if multiline}
	<textarea
		class="rounded-md border border-gray-300 p-1.5 outline-none focus:border-blue-500 focus:ring-1 focus:ring-blue-500"
		bind:value
		{id}
		{name}
		{placeholder}
		{disabled}
		{rows}
	></textarea>
{:else}
	<input
		type="text"
		class="rounded-md border border-gray-300 p-1.5 outline-none focus:border-blue-500 focus:ring-1 focus:ring-blue-500"
		bind:value
		{id}
		{name}
		{placeholder}
		{disabled}
	/>
{/if}
