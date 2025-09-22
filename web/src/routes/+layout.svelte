<script lang="ts">

	//export const prerender = true;
	import favicon from "$lib/assets/favicon.svg";
	import Header from "../components/Header.svelte";

	let { children } = $props();
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
</svelte:head>

<Header />

<div class="container mt-3 mb-4 mx-auto">
	{@render children?.()}
</div>
