<script lang="ts">
	import { DatePicker } from 'bits-ui';
	import { getLocalTimeZone, fromDate, type DateValue } from '@internationalized/date';
	import CalendarIcon from '@lucide/svelte/icons/calendar';
	import ChevronLeft from '@lucide/svelte/icons/chevron-left';
	import ChevronRight from '@lucide/svelte/icons/chevron-right';
	import type { Snippet } from 'svelte';

	type Props = {
		id: string;
		name: string;
		value?: Date;
		label: string;
	};

	let { id, name, value = $bindable(), label }: Props = $props();

	function getValue() {
		if (value) {
			return fromDate(value, getLocalTimeZone());
		}
	}

	function setValue(newValue: DateValue | undefined) {
		if (newValue) {
			value = newValue.toDate(getLocalTimeZone());
		}
	}
</script>

<DatePicker.Root
	weekdayFormat="short"
	fixedWeeks={true}
	locale="fr-CH"
	granularity="day"
	bind:value={getValue, setValue}
>
	<div class="flex w-full flex-col gap-1.5">
		<DatePicker.Label class="block select-none">
			{label}
		</DatePicker.Label>
		<DatePicker.Input
			{id}
			{name}
			class={[
				'h-input border-border-input text-foreground flex w-full items-center rounded-md border border-gray-300 bg-white px-2 py-1 text-sm select-none',
				' focus-within:border-blue-500 focus-within:ring-1 focus-within:ring-blue-500'
			]}
		>
			{#snippet children({ segments })}
				{#each segments as { part, value }, i (part + i)}
					<div class="inline-block select-none">
						{#if part === 'literal'}
							<DatePicker.Segment {part} class="text-gray-800">
								{value}
							</DatePicker.Segment>
						{:else}
							<DatePicker.Segment
								{part}
								class="focus:text-foreground rounded-md px-1 py-1.5 outline-none hover:bg-gray-100 focus:bg-gray-100 focus-visible:ring-0! focus-visible:ring-offset-0! aria-[valuetext=Empty]:text-gray-800"
							>
								{value}
							</DatePicker.Segment>
						{/if}
					</div>
				{/each}
				<DatePicker.Trigger
					class="ml-auto inline-flex size-8 items-center justify-center rounded-md text-gray-600 transition-all outline-none hover:bg-gray-100 focus:bg-gray-100 active:bg-gray-100"
				>
					<CalendarIcon class="size-6" />
				</DatePicker.Trigger>
			{/snippet}
		</DatePicker.Input>

		<DatePicker.Content sideOffset={6} class="z-50" align="end">
			<DatePicker.Calendar class="shadow-popover rounded-md border border-gray-300 bg-white p-2">
				{#snippet children({ months, weekdays })}
					<DatePicker.Header class="flex items-center justify-between">
						<DatePicker.PrevButton
							class="inline-flex size-10 items-center justify-center rounded-md bg-white transition-all hover:bg-gray-100 active:scale-[0.98]"
						>
							<ChevronLeft class="size-6" />
						</DatePicker.PrevButton>
						<DatePicker.Heading class="font-medium" />
						<DatePicker.NextButton
							class="inline-flex size-10 items-center justify-center rounded-md bg-white transition-all hover:bg-gray-100 active:scale-[0.98]"
						>
							<ChevronRight class="size-6" />
						</DatePicker.NextButton>
					</DatePicker.Header>
					<div class="flex flex-col space-y-4 pt-4 sm:flex-row sm:space-y-0 sm:space-x-4">
						{#each months as month (month.value)}
							<DatePicker.Grid class="w-full border-collapse space-y-1 select-none">
								<DatePicker.GridHead>
									<DatePicker.GridRow class="mb-1 flex w-full justify-between">
										{#each weekdays as day (day)}
											<DatePicker.HeadCell
												class="w-10 rounded-md text-xs font-normal! text-gray-800"
											>
												<div>{day.slice(0, 2)}</div>
											</DatePicker.HeadCell>
										{/each}
									</DatePicker.GridRow>
								</DatePicker.GridHead>
								<DatePicker.GridBody>
									{#each month.weeks as weekDates (weekDates)}
										<DatePicker.GridRow class="flex w-full">
											{#each weekDates as date (date)}
												<DatePicker.Cell
													{date}
													month={month.value}
													class="relative size-10 p-0! text-center text-sm"
												>
													<DatePicker.Day
														class={[
															'hover:bg-blue-200 data-disabled:text-gray-500 data-selected:bg-blue-500 data-selected:text-white data-unavailable:text-gray-800',
															'data-disabled:pointer-events-none data-outside-month:pointer-events-none data-selected:font-medium data-unavailable:line-through',
															'group relative inline-flex size-10 items-center justify-center rounded-md border border-transparent bg-transparent p-0 text-sm font-normal whitespace-nowrap transition-all hover:border-white'
														]}
													>
														<div
															class="absolute top-[5px] hidden size-1 rounded-full bg-white transition-all group-data-selected:bg-blue-300 group-data-today:block"
														></div>
														{date.day}
													</DatePicker.Day>
												</DatePicker.Cell>
											{/each}
										</DatePicker.GridRow>
									{/each}
								</DatePicker.GridBody>
							</DatePicker.Grid>
						{/each}
					</div>
				{/snippet}
			</DatePicker.Calendar>
		</DatePicker.Content>
	</div>
</DatePicker.Root>
