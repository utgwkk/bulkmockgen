# CHANGELOG

## Next

- Rename mockgengen to **bulkmockgen**

## Version 0.1.0 (2023/7/8)

- **incompatible**: switch to wrap mockgen command
  - You can use mockgengen with mockgen's command line options.
  - eg. `mockgengen MockBars -- -package mock_bar -destination ./mock_bar/mock_bar.go`

## Version 0.0.2 (2023/7/8)

- migrator: treat consecutive go:generate comment correctly

## Version 0.0.1 (2023/7/8)

- First release
