<script lang="ts">
	/**   */

	import { highlightCode } from './prism';

	let { code = '', language = 'typescript', className = '', children = () => null } = $props();

	// Normalize code formatting
	function normalizeCode(inputCode: string): string {
		const lines = inputCode.split('\n');
		const tabSize = 4; // Standard tab size

		// Remove leading empty lines
		while (lines.length > 0 && lines[0].trim() === '') {
			lines.shift();
		}

		// Remove trailing empty lines
		while (lines.length > 0 && lines[lines.length - 1].trim() === '') {
			lines.pop();
		}

		// Convert tabs to spaces and find minimum indentation
		let minIndent = Infinity;
		const expandedLines = lines.map((line) => {
			// Skip empty lines for min indent calculation
			if (line.trim() === '') return line;

			// Expand tabs to spaces
			let expanded = '';
			let column = 0;

			for (let j = 0; j < line.length; j++) {
				if (line[j] === '\t') {
					// Calculate spaces needed to reach next tab stop
					const spacesToAdd = tabSize - (column % tabSize);
					expanded += ' '.repeat(spacesToAdd);
					column += spacesToAdd;
				} else {
					expanded += line[j];
					column++;
				}
			}

			// Count leading spaces in expanded line
			let indent = 0;
			for (let j = 0; j < expanded.length; j++) {
				if (expanded[j] === ' ') {
					indent++;
				} else {
					break;
				}
			}
			minIndent = Math.min(minIndent, indent);

			return expanded;
		});

		// If no non-empty lines found, return empty string
		if (minIndent === Infinity) return '';

		// Remove common indentation from all lines
		if (minIndent > 0) {
			const normalized = expandedLines.map((line) => {
				// For empty lines, return empty string
				if (line.trim() === '') return '';
				// For non-empty lines, remove the common indentation
				const result = line.substring(minIndent);
				return result;
			});
			return normalized.join('\n');
		}

		return expandedLines.join('\n');
	}

	// For single language mode
	let normalizedCode = $derived(code ? normalizeCode(code) : '');
	let highlightedCode = $derived(normalizedCode ? highlightCode(normalizedCode, language) : '');

	// For multi-language mode, we'll process each CodeLang component
	let codeElement: HTMLElement | undefined = $state();

	function processSlots() {
		if (!codeElement) return;

		const langSections = codeElement.querySelectorAll('[data-lang]');
		if (langSections.length > 0) {
			langSections.forEach((section) => {
				const lang = section.getAttribute('data-lang') || 'text';
				const content = section.textContent || '';
				const normalizedContent = normalizeCode(content);
				const highlighted = highlightCode(normalizedContent, lang);
				section.innerHTML = highlighted;
			});
		}
	}

	$effect(() => {
		if (codeElement) {
			processSlots();
		}
	});
</script>

<pre class="overflow-x-auto bg-gray-100 p-6 {className}">{#if code}<code
			bind:this={codeElement}
			class="language-{language}"
			><!-- eslint-disable svelte/no-at-html-tags -->{@html highlightedCode}</code
		>{:else}{@render children?.()}{/if}</pre>
