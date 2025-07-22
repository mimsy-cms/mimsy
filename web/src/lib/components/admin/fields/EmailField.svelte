<script lang="ts">
	type Props = {
		id: string;
		name: string;
		value: string;
		placeholder?: string;
		disabled?: boolean;
	};

	let { id, name, value = $bindable(), placeholder, disabled }: Props = $props();
</script>

<input
	type="email"
	class="rounded-md border border-gray-300 p-1.5 outline-none focus:border-blue-500 focus:ring-1 focus:ring-blue-500"
	bind:value
	{id}
	{name}
	{placeholder}
	{disabled}
/>
