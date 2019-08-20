// The terror package implements typed, user friendly errors in go. The typed errors parallel and are compatible with HTTP errors. A typed error has two basic components:
// * A user friendly error message meant for viewing by an end-user.
// * The original underlying error, if any.
//
// Note that terror is meant to augment, not replace, basic go errors. In particular, terror is not very useful when dealing with pure backend libraries or anyplace where all the errors would by necessity be 'ServerError' types. Rather, terror is meant to be used primarily at the boundary of user input and request handling where there's the possibility that the problem is something other than a programming bug or backend systems problem.
//
// Another possible use case would be in an unrelable network environment where 'retries' would be a common occurance. In general, terror is useful when the caller might take different actions based on the error type and not so useful when all the caller can do is give up gracefully.
//
// Implementations can thus easily support both communicating directly to the user while maintaining all the raw information one would want to capture in a log.
//
// terror also supports easy debug logging by either calling 'EchoErrorLog' from your test file or setting the `DEBUG_TERROR` environment variable to any non-blank value.
package terror
