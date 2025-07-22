<script lang="ts">
	import { cn } from '$lib/cn';

	type Props = {
		href: string;
		url: string;
		alt: string;
		class?: string;
	};

	let { url, href, alt, class: className }: Props = $props();
</script>

<a class={cn('overflow-hidden rounded-md bg-gray-200', className)} {href}>
	<img src={url} {alt} class="object-cover" />
</a>
