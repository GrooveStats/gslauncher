# GrooveStats Launcher

## What is it

The launcher allows StepMania themes to connect to
[GrooveStats](https://www.groovestats.com/) via it's API. A separate program is
necessary for that because StepMania themes can not make network requests.

- [Screenshots!]
- For pad players only
- Has the following features
- Auto score submission
- Leaderboards accessible via the menu
- Top scores in the songwheel
- SRPG integration, including auto-downloading unlocks
- Hidden dogecoin miner


## How to Install

- Install StepMania. Currently supported are versions 5.0, 5.1 beta 2 and 5.3
  (also know as Outfox. Please use the latest alpha release).
- Set up the Simply Love theme.
  [Installation](https://github.com/Simply-Love/Simply-Love-SM5#installing-simply-love).
- Note: For now, you need to use the `beta` or `gs` branches for this to work.
- Install the launcher. You can get executable files for Windows and Linux
  [here](https://github.com/GrooveStats/gslauncher/releases). Alternatively you
  can also
  [build it yourself](https://github.com/GrooveStats/gslauncher/blob/main/doc/building.md).
  Check out the settings after the first laucnh and adapt the path to your
  StepMania executable and to the data directory.
- Enable local profiles if you haven't already. Alternatively, you can use USB profiles.
- You'll need to generate a `GrooveStats.ini` file in your profile folder. This can be done in one of two ways:
  - Enter the Music Select screen when logged into a profile. This will automotically generate the file for you.
  - Identify your LocalProfile, and manually create a `GrooveStats.ini` with the following contents:

  ```
  [GrooveStats]
  ApiKey=YOUR_API_KEY
  IsPadPlayer=1
  ```

  - IsPadPlayer=1 indicates that you are playing on a pad (and not a keyboard).
    If you're using a keyboard, then set this to 0. In that case your scores
    will not be submitted to GrooveStats, but other functionality will be
    active.

- To obtain your API key to make requests to the GrooveStats service, log into
  your GrooveStats account and visit the
  [Update Profile](https://groovestats.com/index.php?page=register&action=update)
  page to generate and copy your API key. Paste the API key after the `ApiKey=`
  row in the GrooveStats.ini file.
 
- Still have questions or run into problems? Visit the [GrooveStats Discord](https://discord.gg/H7jYZ7xaEX) and ask for help.


## Links for Developers & Themers
- [Building the Groovestats Launcher](https://github.com/GrooveStats/gslauncher/blob/main/doc/building.md)
- [Filesystem IPC](https://github.com/GrooveStats/gslauncher/blob/main/doc/fsipc.md):
  The filesystem based protocol used to communicate between the Theme/StepMania
  and the GrooveStats Launcher.

Are you a theme developer who is interested in integrating the gslauncher into
your theme? Please [contact us](https://discord.gg/H7jYZ7xaEX) before doing so!
