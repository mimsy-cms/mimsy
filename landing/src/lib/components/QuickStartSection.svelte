<script lang="ts">
	import CodeHighlight from './code/CodeHighlight.svelte';
	import CodeLang from './code/CodeLang.svelte';
	import CodeComment from './code/CodeComment.svelte';
</script>

<section class="container mx-auto border-t border-gray-200 px-4 py-16">
	<div class="mx-auto max-w-4xl">
		<h2 class="mb-8 text-3xl font-bold md:text-4xl">Get started in minutes</h2>

		<CodeHighlight className="mb-8">
			<CodeComment># 1. Install mimsy CLI</CodeComment>
			<CodeLang lang="bash" code="npm install -g mimsy-cli" />

			<CodeComment># 2. Initialize project</CodeComment>
			<CodeLang lang="bash" code="mimsy init my-cms" />

			<CodeComment># 3. Define your first collection</CodeComment>
			<CodeComment># collections/posts.ts</CodeComment>
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
			<CodeComment># 4. Build and deploy</CodeComment>
			<CodeLang
				lang="bash"
				code={`mimsy build
mimsy deploy`}
			/>
		</CodeHighlight>

		<div class="text-center">
			<a
				href="https://github.com/[your-org]/mimsy"
				class="inline-block bg-black px-6 py-3 text-white transition-colors hover:bg-gray-800"
			>
				View Full Documentation on GitHub →
			</a>
		</div>
	</div>
</section>
