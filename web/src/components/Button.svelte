<script>
  export let label = "Click Me";

  function handleClick() {
    alert("Button clicked!");
  }

</script>

<style>
  button {
    background-color: #03dac5;
    border: none;
    color: white;
    padding: 0.5rem 1rem;
    cursor: pointer;
  }
</style>

<button on:click={handleClick}>{label}</button>