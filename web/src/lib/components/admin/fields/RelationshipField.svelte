<script lang="ts">
	import type { FieldRelation } from '$lib/collection/definition';
	import SelectField from './SelectField.svelte';

	type Props = {
		name: string;
		value?: string;
		label: string;
		field: FieldRelation;
	};

	type Item = {
		value: string;
		label: string;
	};

	function getEndpoint(resourceSlug: string) {
		switch (resourceSlug) {
			case '<builtins.user>':
				return '/api/v1/users';
			case '<builtins.media>':
				return '/api/v1/media';
			default:
				return `/api/v1/collections/${resourceSlug}`;
		}
	}

	let { name, label, field, value = $bindable() }: Props = $props();

	let items = $state<Item[]>([]);

	async function fetchResources() {
		const response = await fetch(getEndpoint(field.relatesTo));
		const resources = await response.json();

		switch (field.relatesTo) {
			case '<builtins.user>':
				items = resources.map((r: any) => ({ value: r.id, label: r.email }));
				break;
			case '<builtins.media>':
				items = resources.map((r: any) => ({ value: r.id, label: r.name }));
				break;
			default:
				items = resources.map((r: any) => ({ value: r.id, label: r.slug }));
		}
	}

	$effect(() => {
		fetchResources();
	});
</script>

<SelectField {name} {label} {items} bind:value />
