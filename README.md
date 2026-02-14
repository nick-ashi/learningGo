# LEARNING GO

This repository is for me and anyone else who finds it to use as a resource for learning Go. 

## WHAT'S IN THE REPO?

There's one file with a bunch of random notes on various aspects of GoLang -- environment setup, importing basic packages like fmt, variable declarations, types, implicit and explicit assignments, for loops, console output, etc. etc. Kind of liking making my own learn x in y document!

Currently the first little project is just a cli todo application for me to learn the fundamentals of the language while throwing in some read/write ops for JSON files. Also of course since it's using the cli it uses a couple os operations as well. Future projects may or may not come -- depends on whether or not I have / how much time I commit to diving in here.

### PROJECTS

**TODO CLI**
Simple todo cli application to get familiar with GoLang

**POKEDEX CLI**
Similar to pokefetch. Used claude to do a lot of the rendering of non-pokemon-colorscript supported sprites.
- Fetches info from PokeAPI, sprites from Pokemon-Colorscripts
- If we don't get a 200 stat code from the pokemon-colorscripts repo, we fallback to rendering a pokeapi-provided sprite ourselves
- Prints the info on the righthand side of the sprite in a very pokefetch like fashion
- shiny flag option available
- PROBABLY going to implement some option to do only the sprite


## WHY LEARN GO?

Honestly, it's my own curiosity... I've heard it's a fun language to learn and use, and being a compiled language created with the intention of being a kind of bridge between low-level langs like C++ and the easy to understand languages like Python. Also hoping to dive deeper into topics like concurrency, multithreading, etc.

## Future Objectives / Learning Goals
- [ ] Want to learn more in-depth backend web development with Go (i am getting tired of using javascript)
  - [ ] learn how to implement APIs with Go
- [ ] Build some CLI applications
  - [ ] Pokedex CLI -- WIP 
- [ ] Learn about threading / concurrency with Go
- More ???