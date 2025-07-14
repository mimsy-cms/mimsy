<script lang="ts">
	import CodeHighlight from './code/CodeHighlight.svelte';
	import CodeLang from './code/CodeLang.svelte';
	import CodeComment from './code/CodeComment.svelte';
</script>

<section class="container mx-auto border-t border-gray-200 px-4 py-16">
	<div class="mx-auto max-w-4xl">
		<h2 class="mb-8 text-3xl font-bold md:text-4xl">Technical Overview</h2>

		<pre class="mb-12 overflow-x-auto text-xs md:text-sm">
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Developer     │    │   mimsy CLI     │    │   Go Backend    │
│   (TypeScript)  │───▶│   (Schema)      │───▶│   (REST API)    │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                                                      │
                                                      ▼
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Admin Panel   │    │   Your App      │    │   PostgreSQL    │
│   (SvelteKit)   │◀───│   (SvelteKit)   │◀───│   + pgroll      │
└─────────────────┘    └─────────────────┘    └─────────────────┘
		</pre>

		<h3 class="mb-4 text-2xl font-bold">Workflow Example</h3>

		<CodeHighlight>
			<CodeComment>// 1. Define your schema in TypeScript</CodeComment>
			<CodeComment>// collections/blog.ts</CodeComment>
			<CodeLang
				lang="typescript"
				code={`
				export const BlogPost = collection({
                  name: "post",
                  fields: [
                    { name: "title", type: "text", required: true },
                    { name: "content", type: "richtext" },
                    { name: "publishedDate", type: "date" },
                  ],
                });`}
			/>
			<CodeComment>// 2. Process with CLI</CodeComment>
			<CodeLang lang="bash" code="$ mimsy build" />
			<CodeComment>// 3. Deploy with zero downtime</CodeComment>
			<CodeLang lang="bash" code="$ mimsy deploy" />
			<CodeLang
				lang="text"
				code={`✓ Schema migration created
✓ Old API version: active
✓ New API version: active
✓ Traffic migration: complete`}
			/>
		</CodeHighlight>
	</div>
</section>
