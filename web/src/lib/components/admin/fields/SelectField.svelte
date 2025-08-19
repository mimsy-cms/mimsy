<script lang="ts">
	import { Select } from 'bits-ui';
	import CheckIcon from '@lucide/svelte/icons/check';
	import ChevronsUpDownIcon from '@lucide/svelte/icons/chevrons-up-down';

	type Item = {
		value: string;
		label: string;
	};

	type Props = {
		name: string;
		value?: Props['multiple'] extends true ? string[] : string;
		label: string;
		items: Item[];
		multiple?: boolean;
	};

	let { name, value = $bindable(), label, items, multiple }: Props = $props();

	const selectedLabel = $derived.by(() => {
		if (!value || (Array.isArray(value) && value.length === 0)) return label;

		if (multiple) {
			return items
				.filter((item) => Array.isArray(value) && value.includes(item.value))
				.map((item) => item.label)
				.join(', ');
		}
		return items.find((item) => item.value === value)?.label;
	});
</script>

<Select.Root
	type={multiple ? 'multiple' : 'single'}
	onValueChange={(v: any) => (value = v)}
	{name}
	{items}
>
	<Select.Trigger
		class={[
			'inline-flex touch-none items-center rounded-md border border-gray-300 bg-white py-2 text-sm select-none data-placeholder:text-gray-600',
			'data-[state=open]:border-blue-500 data-[state=open]:ring-1 data-[state=open]:ring-blue-500'
		]}
		aria-label={label}
	>
		<span class="ml-3">{selectedLabel}</span>
		<ChevronsUpDownIcon class="mr-2 ml-auto size-4 text-gray-600" />
	</Select.Trigger>
	<Select.Portal>
		<Select.Content
			class={[
				'data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 data-[state=closed]:zoom-out-99',
				'data-[state=open]:zoom-in-99 data-[side=bottom]:slide-in-from-top-1 data-[side=left]:slide-in-from-right-1 data-[side=right]:slide-in-from-left-1 data-[side=top]:slide-in-from-bottom-1',
				'max-h-[var(--bits-select-content-available-height)] w-[var(--bits-select-anchor-width)] min-w-[var(--bits-select-anchor-width)] shadow-md outline-hidden',
				'rounded-xl border border-gray-300 bg-white px-1 py-2 select-none data-[side=bottom]:translate-y-1 data-[side=left]:-translate-x-1 data-[side=right]:translate-x-1 data-[side=top]:-translate-y-1'
			]}
			side="bottom"
			sideOffset={10}
		>
			<Select.Viewport class="flex flex-col gap-1 p-1">
				{#each items as item, i (i + item.value)}
					<Select.Item
						class={[
							'flex  h-10 w-full items-center rounded-md px-3 py-1.5 text-sm capitalize outline-hidden select-none data-disabled:opacity-50 data-highlighted:bg-gray-100',
							'data-selected:bg-blue-500 data-selected:text-white'
						]}
						value={item.value}
						label={item.label}
					>
						{#snippet children({ selected })}
							{item.label}
							{#if selected}
								<div class="ml-auto">
									<CheckIcon class="size-4 stroke-[2.5]" aria-label="check" />
								</div>
							{/if}
						{/snippet}
					</Select.Item>
				{/each}
			</Select.Viewport>
		</Select.Content>
	</Select.Portal>
</Select.Root>
