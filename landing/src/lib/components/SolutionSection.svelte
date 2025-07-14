<script lang="ts">
	import CodeHighlight from './code/CodeHighlight.svelte';
</script>

<section class="container mx-auto border-t border-gray-200 px-4 py-16">
	<div class="mx-auto max-w-4xl">
		<h2 class="mb-8 text-3xl font-bold md:text-4xl">How mimsy solves this</h2>

		<div class="space-y-8">
			<div>
				<h3 class="mb-2 text-xl font-bold">1. CODE IS TRUTH</h3>
				<div class="space-y-1 pl-6 text-gray-700">
					<div>Schema lives in TypeScript files</div>
					<CodeHighlight
						className="ml-5"
						language="typescript"
						code={`
						collection({
                          name: "post",
                          fields: [
                            { name: "title", type: "text", required: true },
                            { name: "content", type: "richtext" },
                            { name: "publishedDate", type: "date" },
                          ],
                        });
					`}
					/>
					<div>No UI modifications allowed, pure code only</div>
					<div>Git tracks all changes, and allows for versioned access</div>
				</div>
			</div>

			<div>
				<h3 class="mb-2 text-xl font-bold">2. ZERO DOWNTIME</h3>
				<div class="space-y-1 pl-6 text-gray-700">
					<div>pgroll-powered live migrations</div>
					<div>Progressive rollouts made easy</div>
					<div>Rollback without service interruption</div>
				</div>
			</div>

			<div>
				<h3 class="mb-2 text-xl font-bold">3. TYPESCRIPT NATIVE</h3>
				<div class="space-y-1 pl-6 text-gray-700">
					<div>Type-safe resource definitions</div>
					<div>Compile-time validation</div>
					<div>Build your schema in typescript, get an SDK for free</div>
				</div>
			</div>

			<div>
				<h3 class="mb-2 text-xl font-bold">4. SVELTEKIT FIRST</h3>
				<div class="space-y-1 pl-6 text-gray-700">
					<div>Built by a SvelteKit team for other SvelteKit teams</div>
					<div>Only SvelteKit, no auto-generated tooling for other frameworks</div>
					<div>Use what you already know, enjoy a simple learning curve</div>
				</div>
			</div>

			<div>
				<h3 class="mb-2 text-xl font-bold">5. ABSOLUTELY. NO. BULLSHIT.</h3>
				<div class="space-y-1 pl-6 text-gray-700">
					<div>We are tired of products showing off fancy but useless features for years</div>
					<div>So we made sure to not overpromise or show things that aren't useful</div>
					<div>Just the bare minimum to get working, no more, no less</div>
					<div>
						This website serves as our manifesto, no fancy pictures or animations to see around here
					</div>
				</div>
			</div>
		</div>
	</div>
</section>
