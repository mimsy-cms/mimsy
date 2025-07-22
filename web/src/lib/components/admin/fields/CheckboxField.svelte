<script lang="ts">
	import { Checkbox, Label } from 'bits-ui';
	import CheckIcon from '@lucide/svelte/icons/check';
	import { cn } from '$lib/cn';

	type Props = {
		id: string;
		name: string;
		checked?: boolean;
		label: string;
		class?: string;
	};

	let { id, name, checked = $bindable(), label, class: className }: Props = $props();
</script>

<div class={cn('flex items-center space-x-3', className)}>
	<Checkbox.Root
		{id}
		{name}
		aria-labelledby={`${id}-label`}
		class={[
			'border-blue-900 bg-blue-500 data-[state=unchecked]:border-gray-300 data-[state=unchecked]:bg-white',
			'peer inline-flex size-6 items-center justify-center rounded-md border transition-all duration-75 ease-in-out active:scale-[0.98]'
		]}
	>
		{#snippet children({ checked })}
			<div class="inline-flex items-center justify-center text-white">
				{#if checked}
					<CheckIcon class="size-4 stroke-[2.5]" />
				{/if}
			</div>
		{/snippet}
	</Checkbox.Root>
	<Label.Root id={`${id}-label`} for={id} class="leading-none peer-disabled:cursor-not-allowed">
		{label}
	</Label.Root>
</div>
