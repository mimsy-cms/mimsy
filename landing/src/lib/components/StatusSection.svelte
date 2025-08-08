<script lang="ts">
	import { onMount, onDestroy } from 'svelte';

	let spinnerIndex = 0;
	let spinnerInterval: number | NodeJS.Timeout;
	let progressBars: { [key: string]: number } = {
		typescript: 35,
		sveltekit: 60,
		pgroll: 25
	};

	const spinnerChars = ['⠋', '⠙', '⠹', '⠸', '⠼', '⠴', '⠦', '⠧', '⠇', '⠏'];

	onMount(() => {
		spinnerInterval = setInterval(() => {
			spinnerIndex = (spinnerIndex + 1) % spinnerChars.length;
		}, 100);
	});

	onDestroy(() => {
		if (spinnerInterval) {
			clearInterval(spinnerInterval);
		}
	});

	function generateProgressBar(percentage: number): string {
		const width = 20;
		const filled = Math.floor((percentage / 100) * width);
		const empty = width - filled;
		return '▓'.repeat(filled) + '░'.repeat(empty);
	}
</script>

<section class="container mx-auto border-t border-gray-200 px-4 py-16">
	<div class="mx-auto max-w-4xl">
		<h2 class="mb-8 text-3xl font-bold md:text-4xl">Status & Roadmap</h2>

		<div class="bg-gray-100 p-6">
			<div class="mb-8">
				<div class="cli-line">
					<span class="font-bold text-green-700">></span> mimsy status --verbose
				</div>
				<div class="cli-line">
					STATUS: <span class="text-yellow-700">Early Development</span>
				</div>
				<div class="cli-line">
					VERSION: <span class="text-cyan-700">0.1.0</span>
				</div>
				<div class="cli-line">
					STAGE: <span class="text-orange-700">Architecture & Core Implementation</span>
				</div>
			</div>

			<div>
				<div class="cli-line">
					<span class="font-bold text-green-700">></span> mimsy roadmap --progress
				</div>
				<div class="cli-line">
					<span class="font-bold text-green-700">✓</span> Core architecture design
				</div>
				<div class="cli-line">
					<span class="font-bold text-green-700">✓</span> Go backend foundation
				</div>
				<div class="cli-line">
					<span class="spinner font-bold text-yellow-700">{spinnerChars[spinnerIndex]}</span>
					TypeScript SDK development
					<span class="ml-auto text-xs text-cyan-700"
						>[{generateProgressBar(progressBars.typescript)}] {progressBars.typescript}%</span
					>
				</div>
				<div class="cli-line">
					<span class="spinner font-bold text-yellow-700">{spinnerChars[spinnerIndex]}</span>
					SvelteKit admin panel
					<span class="ml-auto text-xs text-cyan-700"
						>[{generateProgressBar(progressBars.sveltekit)}] {progressBars.sveltekit}%</span
					>
				</div>
				<div class="cli-line">
					<span class="spinner font-bold text-yellow-700">{spinnerChars[spinnerIndex]}</span> pgroll
					integration
					<span class="ml-auto text-xs text-cyan-700"
						>[{generateProgressBar(progressBars.pgroll)}] {progressBars.pgroll}%</span
					>
				</div>
				<div class="cli-line">
					<span class="text-gray-500">◯</span> Beta release
				</div>
				<div class="cli-line">
					<span class="text-gray-500">◯</span> Production ready
				</div>
			</div>
		</div>
	</div>
</section>

<style>
	.cli-line {
		margin-bottom: 4px;
		display: flex;
		align-items: center;
		gap: 8px;
		font-family:
			'Fira Code', 'SF Mono', Monaco, 'Cascadia Code', 'Roboto Mono', Consolas, 'Courier New',
			monospace;
		font-size: 14px;
		line-height: 1.6;
	}

	.spinner {
		animation: pulse 1s infinite;
	}

	@keyframes pulse {
		0%,
		100% {
			opacity: 1;
		}
		50% {
			opacity: 0.7;
		}
	}
</style>
