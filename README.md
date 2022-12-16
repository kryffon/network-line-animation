# Network Line Animation

Built using [ebitengine.org](https://ebitengine.org/)

![sample](https://user-images.githubusercontent.com/120729657/208117098-b79c154a-5297-4d6f-9ec2-cc30d121e0fc.png)

### Optimizations implemented

- Parallel workers for finding nearest neighbours (no significant improvement)
- [K-d tree](https://en.wikipedia.org/wiki/K-d_tree) for finding nearest neighbours (huge improvement). Used package [kdbush](https://github.com/MadAppGang/kdbush)
