# Network Line Animation

Built using [ebitengine.org](https://ebitengine.org/)


### Optimizations implemented

- Parallel workers for finding nearest neighbours (no significant improvement)
- [K-d tree](https://en.wikipedia.org/wiki/K-d_tree) for finding nearest neighbours (huge improvement). Used package [kdbush](https://github.com/MadAppGang/kdbush)