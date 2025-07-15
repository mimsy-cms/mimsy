<script lang="ts">
	import HeroSection from '$lib/components/HeroSection.svelte';
	import ProblemSection from '$lib/components/ProblemSection.svelte';
	import SolutionSection from '$lib/components/SolutionSection.svelte';
	import TechnicalOverviewSection from '$lib/components/TechnicalOverviewSection.svelte';
	import FeaturesSection from '$lib/components/FeaturesSection.svelte';
	import QuickStartSection from '$lib/components/QuickStartSection.svelte';
	import StatusSection from '$lib/components/StatusSection.svelte';
	import TeamSection from '$lib/components/TeamSection.svelte';
	import FooterSection from '$lib/components/FooterSection.svelte';
</script>

<svelte:head>
	<title>mimsy - A headless CMS for agencies who code</title>
	<meta name="description" content="Zero-downtime deployments. Code-first schema. TypeScript native. Built for SvelteKit teams who want control over their CMS." />
</svelte:head>

<div class="min-h-screen bg-white text-black font-mono">
	<HeroSection />
	<ProblemSection />
	<SolutionSection />
	<TechnicalOverviewSection />
	<FeaturesSection />
	<QuickStartSection />
	<StatusSection />
	<TeamSection />
	<FooterSection />
</div>

<style>
	:global(pre) {
		white-space: pre;
		line-height: 1.2;
	}
	
	:global(code) {
		font-family: 'Fira Code', monospace;
	}
</style>