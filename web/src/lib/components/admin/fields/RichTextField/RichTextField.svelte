<script lang="ts">
import { onMount } from 'svelte';
import { Editor } from '@tiptap/core';
import StarterKit from '@tiptap/starter-kit';

let { value = $bindable(''), ...props } = $props();

let element: HTMLElement;
let editor: Editor;
let isUpdatingFromExternal = false;

onMount(() => {
	editor = new Editor({
		element: element,
		extensions: [
			StarterKit,
		],
		content: parseValue(value),
		onTransaction: () => {
			editor = editor;
		},
		onUpdate: ({ editor }) => {
			if (!isUpdatingFromExternal) {
				const json = editor.getJSON();
				value = json;
			}
		}
	});

	return () => {
		if (editor) {
			editor.destroy();
		}
	};
});

function parseValue(val: any) {
	if (!val) {
		return '<p></p>';
	}
	
	if (typeof val === 'string') {
		return val;
	}
	
	if (Array.isArray(val)) {
		return val.map(text => `<p>${text}</p>`).join('');
	}
	
	if (typeof val === 'object' && val.type) {
		return val;
	}
	
	return '<p></p>';
}

$effect(() => {
	if (editor && value !== undefined) {
		const currentContent = editor.getJSON();
		const parsedValue = parseValue(value);
		
		if (JSON.stringify(currentContent) !== JSON.stringify(parsedValue)) {
			isUpdatingFromExternal = true;
			editor.commands.setContent(parsedValue, false);
			isUpdatingFromExternal = false;
		}
	}
});
</script>

<div bind:this={element} class="richtext-editor prose max-w-none"></div>

<style>
:global(.richtext-editor) {
	border: 1px solid #d1d5db;
	border-radius: 0.375rem;
	padding: 0.75rem;
	min-height: 120px;
}

:global(.richtext-editor:focus-within) {
	outline: 2px solid #3b82f6;
	outline-offset: 2px;
}

:global(.ProseMirror) {
	outline: none;
}
</style>