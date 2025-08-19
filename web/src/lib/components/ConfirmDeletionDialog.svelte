<script lang="ts">
	import { Dialog } from 'bits-ui';
	import { cn } from '$lib/cn';
	import XIcon from '@lucide/svelte/icons/x';

	type Props = {
		open: boolean;
		title?: string;
		description?: string;
		onConfirm: () => void;
		onCancel: () => void;
	};

	let {
		open = $bindable(),
		title = 'Confirm Deletion',
		description,
		onConfirm,
		onCancel
	}: Props = $props();

	function handleOpenChange(newOpen: boolean) {
		open = newOpen;
		if (!newOpen) {
			onCancel();
		}
	}

	function handleConfirm() {
		onConfirm();
		open = false;
	}

	function handleCancel() {
		onCancel();
		open = false;
	}
</script>

<Dialog.Root {open} onOpenChange={handleOpenChange}>
	<Dialog.Portal>
		<Dialog.Overlay
			class="data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 fixed inset-0 z-50 bg-black/40 backdrop-blur-xs"
		/>
		<Dialog.Content
			class={[
				'fixed top-1/2 left-1/2 z-50 w-full max-w-lg -translate-x-1/2 -translate-y-1/2 transform gap-4 border border-gray-200 bg-white p-6 shadow-lg duration-200 sm:rounded-lg',
				'data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 data-[state=closed]:zoom-out-95 data-[state=open]:zoom-in-95'
			]}
		>
			<div class="flex items-start gap-4">
				<div class="flex-1">
					<Dialog.Title class="text-lg font-semibold text-gray-900">
						{title}
					</Dialog.Title>
					<Dialog.Description class="mt-2 text-sm text-gray-500">
						{description}
					</Dialog.Description>
				</div>

				<Dialog.Close
					class="rounded-sm opacity-70 ring-offset-white transition-opacity hover:opacity-100 focus:ring-2 focus:ring-gray-950 focus:ring-offset-2 focus:outline-none disabled:pointer-events-none"
				>
					<XIcon class="size-4" />
					<span class="sr-only">Close</span>
				</Dialog.Close>
			</div>

			<div class="mt-4 flex flex-col-reverse gap-2 sm:flex-row sm:justify-end sm:space-x-2">
				<button type="button" onclick={handleCancel} class="btn btn-outline"> Cancel </button>
				<button
					type="button"
					onclick={handleConfirm}
					class={cn(
						'inline-flex h-10 items-center justify-center rounded-md bg-red-600 px-4 py-2 text-sm font-medium text-white hover:bg-red-700 focus:ring-2 focus:ring-red-950 focus:ring-offset-2 focus:outline-none disabled:cursor-not-allowed disabled:opacity-50'
					)}
				>
					Confirm
				</button>
			</div>
		</Dialog.Content>
	</Dialog.Portal>
</Dialog.Root>
