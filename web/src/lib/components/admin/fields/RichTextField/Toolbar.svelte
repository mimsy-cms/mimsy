<script lang="ts">
	import {
		BoldButton,
		Divider,
		ItalicButton,
		UnderlineButton,
		StrikethroughButton,
		FormatCodeButton,
		DropDownAlign,
		FontFamilyDropDown,
		FontSizeDropDown,
		Toolbar
	} from 'svelte-lexical';
</script>

<Toolbar>
	{#snippet children({ editor, activeEditor, blockType })}
		<FontFamilyDropDown />
		<FontSizeDropDown />
		<Divider />
		<BoldButton />
		<ItalicButton />
		<UnderlineButton />
		<StrikethroughButton />
		<FormatCodeButton />
		<Divider />
		<DropDownAlign />
	{/snippet}
</Toolbar>
